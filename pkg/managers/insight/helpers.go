package insight

import (
	"fmt"
	"strings"
	"time"

	"github.com/soralabs/zen/pkg/db"
	"github.com/soralabs/zen/pkg/id"
	"github.com/soralabs/zen/pkg/stores"
)

/*
Example output:
- User is working on a REST API project in Go (Category: technical_detail, Confidence: 0.95)
- Project deadline is March 20th (Category: date_time, Confidence: 0.90)
- User is experienced with Python and Go (Category: personal_info, Confidence: 0.85)
- User prefers detailed explanations (Category: preference, Confidence: 0.90)
*/
func (im *InsightManager) formatInsightFragments(insights []db.Fragment) string {
	var formatted []string
	for _, insight := range insights {
		confidence := insight.Metadata["confidence"].(float64)
		category := insight.Metadata["category"].(string)

		formatted = append(formatted, formatInsight(insight.ID, insight.Content, category, confidence, insight.ActorID, insight.Actor.Name))
	}
	return strings.Join(formatted, "\n")
}

/*
Example output:
# Session Insights
- ID: 01H9X2Z4N8K5Y6M7 | Content: User is working on a REST API project in Go | Category: technical_detail | Confidence: 0.95
- ID: 01H9X2Z4N8K5Y6M9 | Content: Project deadline is March 20th | Category: date_time | Confidence: 0.90

# User Insights
- ID: 01H9X2Z4N8K5Y6M8 | Content: User is experienced with Python and Go | Category: personal_info | Confidence: 0.85
- ID: 01H9X2Z4N8K5Y6N1 | Content: User prefers detailed explanations | Category: preference | Confidence: 0.90
*/
func (im *InsightManager) formatAllInsights(sessionInsights, actorInsights []db.Fragment) string {
	var sessionBuilder, actorBuilder strings.Builder

	// Separate insights by type
	for _, insight := range sessionInsights {
		category := insight.Metadata["category"].(string)
		confidence := insight.Metadata["confidence"].(float64)
		actorName := insight.Actor.Name

		formattedInsight := formatInsight(insight.ID, insight.Content, category, confidence, insight.ActorID, actorName)
		sessionBuilder.WriteString(formattedInsight + "\n")
	}

	for _, insight := range actorInsights {
		category := insight.Metadata["category"].(string)
		confidence := insight.Metadata["confidence"].(float64)
		actorName := insight.Actor.Name

		formattedInsight := formatInsight(insight.ID, insight.Content, category, confidence, insight.ActorID, actorName)
		actorBuilder.WriteString(formattedInsight + "\n")
	}

	var result strings.Builder
	if sessionBuilder.Len() > 0 {
		result.WriteString("# Session Insights\n")
		result.WriteString(sessionBuilder.String())
	}

	if actorBuilder.Len() > 0 {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString("# Actor Insights\n")
		result.WriteString(actorBuilder.String())
	}

	return result.String()
}

/*
Example output:
[just now] User (ID: 123): Can you help me with this Go code?
[1 minute ago] Assistant (ID: 456): I'd be happy to help! What specific issue are you having?
[5 minutes ago] User (ID: 123): I'm trying to implement error handling in my REST API
[1 hour ago] Assistant (ID: 456): Let's break down the best practices for error handling in Go APIs
[2 hours ago] User (ID: 123): Thanks, that would be really helpful
*/
func formatMessageHistory(messages []db.Fragment) string {
	var messageHistory strings.Builder
	now := time.Now()

	for _, msg := range messages {
		role := fmt.Sprintf("%s (ID: %s)", msg.Actor.Name, msg.ActorID)
		if msg.Actor.Assistant {
			role = fmt.Sprintf("Assistant (ID: %s)", msg.ActorID)
		}

		minutesAgo := int(now.Sub(msg.CreatedAt).Minutes())
		timeAgo := ""
		switch {
		case minutesAgo < 1:
			timeAgo = "just now"
		case minutesAgo == 1:
			timeAgo = "1 minute ago"
		case minutesAgo < 60:
			timeAgo = fmt.Sprintf("%d minutes ago", minutesAgo)
		case minutesAgo < 120:
			timeAgo = "1 hour ago"
		default:
			timeAgo = fmt.Sprintf("%d hours ago", minutesAgo/60)
		}

		messageHistory.WriteString(fmt.Sprintf("- [%s] %s: %s\n", timeAgo, role, msg.Content))
	}

	return messageHistory.String()
}

/*
Example output:
# Example 1
## Session
- [2 hours ago] Actor (ID: 123): Hi! I need some help with a REST API I'm building in Go.
- [2 hours ago] Assistant (ID: 456): I'd be happy to help with your Go REST API.

## Extracted Insights
### Session Insights
- Actor is working on a REST API project in Go (Category: technical_detail, Confidence: 0.95)
- Project deadline is March 20th (Category: date_time, Confidence: 0.90)

### Actor Insights
- Actor is experienced with Python and Go (Category: personal_info, Confidence: 0.85)

---
*/
func formatExampleInsights(examples []ExampleSet) string {
	var result strings.Builder

	for i, example := range examples {
		result.WriteString(fmt.Sprintf("\nEXAMPLE %d:\n", i+1))
		result.WriteString("Conversation:\n")
		for _, msg := range example.Messages {
			result.WriteString(fmt.Sprintf("[%s] %s: %s\n", msg.TimeAgo, msg.Role, msg.Content))
		}

		result.WriteString("\nExtracted Insights:\n")
		var conversationInsights, actorInsights strings.Builder

		for _, insight := range example.Insights {
			formatted := fmt.Sprintf("- %s (Category: %s, Confidence: %.2f)\n",
				insight.Content,
				insight.Category,
				insight.Confidence,
			)

			if insight.Type == string(SessionInsights) {
				conversationInsights.WriteString(formatted)
			} else {
				actorInsights.WriteString(formatted)
			}
		}

		if conversationInsights.Len() > 0 {
			result.WriteString("SESSION INSIGHTS:\n")
			result.WriteString(conversationInsights.String())
		}
		if actorInsights.Len() > 0 {
			result.WriteString("ACTOR INSIGHTS:\n")
			result.WriteString(actorInsights.String())
		}
		result.WriteString("\n---\n")
	}

	return result.String()
}

func getUniqueInsights(similarInsights, sessionInsights, actorInsights []db.Fragment) []db.Fragment {
	// Create a map to track existing insights by content
	existingContent := make(map[string]bool)

	// Add session and actor insights content to the map
	for _, insight := range sessionInsights {
		existingContent[insight.Content] = true
	}
	for _, insight := range actorInsights {
		existingContent[insight.Content] = true
	}

	// Only include similar insights that don't exist in session or actor insights
	var uniqueInsights []db.Fragment
	for _, insight := range similarInsights {
		if !existingContent[insight.Content] {
			uniqueInsights = append(uniqueInsights, insight)
			existingContent[insight.Content] = true // Prevent duplicates within similar insights
		}
	}

	return uniqueInsights
}

func getExampleSets() []ExampleSet {
	return []ExampleSet{
		{
			Messages: []ExampleMessage{
				{Role: "Actor (ID: 123)", Content: "Hi! I need some help with a REST API I'm building in Go. I'm trying to set up proper database migrations.", TimeAgo: "2 hours ago"},
				{Role: "Assistant (ID: 456)", Content: "I'd be happy to help with your Go REST API and database migrations. What specific challenges are you facing?", TimeAgo: "2 hours ago"},
				{Role: "Actor (ID: 123)", Content: "Well, I need to get this done by March 20th, and I'm not sure about the best practices for handling migrations in Go.", TimeAgo: "2 hours ago"},
				{Role: "Assistant (ID: 456)", Content: "I understand you're working with a deadline. Let me show you some examples of how to handle database migrations in Go...", TimeAgo: "2 hours ago"},
			},
			Insights: []ExampleInsight{
				{
					Content:    "Actor is working on a REST API project in Go",
					Category:   "technical_detail",
					Confidence: 0.95,
					Type:       string(SessionInsights),
					ActorID:    id.ID("123"),
				},
				{
					Content:    "The REST API project deadline is March 20th",
					Category:   "date_time",
					Confidence: 0.90,
					Type:       string(SessionInsights),
					ActorID:    id.ID("123"),
				},
				{
					Content:    "Actor needs help with database migrations",
					Category:   "requirement",
					Confidence: 0.85,
					Type:       string(SessionInsights),
					ActorID:    id.ID("123"),
				},
			},
		},
		{
			Messages: []ExampleMessage{
				{Role: "Actor (ID: 123)", Content: "Could you help me understand this Python code? I usually work with Go and Python, but this particular pattern is new to me.", TimeAgo: "30 minutes ago"},
				{Role: "Assistant (ID: 456)", Content: "Of course! Since you're familiar with both Go and Python, I'll explain with detailed examples comparing both languages.", TimeAgo: "29 minutes ago"},
				{Role: "Actor (ID: 123)", Content: "Thanks! That's exactly what I need. I'm in EST timezone by the way, so I'll be around for the next few hours if we need to discuss more.", TimeAgo: "28 minutes ago"},
			},
			Insights: []ExampleInsight{
				{
					Content:    "Actor is experienced with Python and Go",
					Category:   "personal_info",
					Confidence: 0.95,
					Type:       string(ActorInsights),
					ActorID:    id.ID("123"),
				},
				{
					Content:    "Actor prefers detailed explanations with code examples",
					Category:   "preference",
					Confidence: 0.90,
					Type:       string(ActorInsights),
					ActorID:    id.ID("123"),
				},
				{
					Content:    "Actor works in EST timezone",
					Category:   "personal_info",
					Confidence: 0.85,
					Type:       string(ActorInsights),
					ActorID:    id.ID("123"),
				},
			},
		},
	}
}

// getInsightData retrieves all relevant insights for a given message fragment.
func (im *InsightManager) getInsightData(message *db.Fragment) (*cachedInsightData, error) {
	data := &cachedInsightData{
		SessionInsights: []db.Fragment{},
		ActorInsights:   []db.Fragment{},
		SimilarInsights: []db.Fragment{},
	}

	// Get session insights
	sessionInsights, err := im.FragmentStore.SearchByFilter(stores.FragmentFilter{
		SessionID: &message.SessionID,
		Metadata: []stores.MetadataCondition{
			{
				Key:      "type",
				Value:    string(SessionInsights),
				Operator: stores.MetadataOpEquals,
			},
		},
		Limit: 10,
	})
	if err != nil {
		return data, fmt.Errorf("failed to get session insights: %w", err)
	}
	data.SessionInsights = append(data.SessionInsights, sessionInsights...)

	// Get actor insights
	actorInsights, err := im.FragmentStore.SearchByFilter(stores.FragmentFilter{
		ActorID: &message.ActorID,
		Metadata: []stores.MetadataCondition{
			{
				Key:      "type",
				Value:    string(ActorInsights),
				Operator: stores.MetadataOpEquals,
			},
		},
		Limit: 10,
	})
	if err != nil {
		return data, fmt.Errorf("failed to get actor insights: %w", err)
	}
	data.ActorInsights = append(data.ActorInsights, actorInsights...)

	// Get similar insights if message has an embedding
	if len(message.Embedding.Slice()) > 0 {
		similarInsights, err := im.FragmentStore.SearchSimilar(message.Embedding, message.SessionID, 3)
		if err != nil {
			return data, fmt.Errorf("failed to get similar insights: %w", err)
		}
		data.SimilarInsights = append(data.SimilarInsights, similarInsights...)
	}

	return data, nil
}

// formatInsight formats a single insight with its metadata into a string
// It provides a consistent format for displaying insight information
func formatInsight(id id.ID, content string, category string, confidence float64, actorID id.ID, actorName string) string {
	return fmt.Sprintf("- %s (ID: %s, ActorID: %s, Actor: %s, Category: %s, Confidence: %.2f)",
		content,
		id,
		actorID,
		actorName,
		category,
		confidence,
	)
}

// removeInsightByID removes an insight from a slice by its ID
// Used when insights become outdated or irrelevant
func removeInsightByID(insights []db.Fragment, idToRemove id.ID) []db.Fragment {
	filtered := make([]db.Fragment, 0, len(insights))
	for _, insight := range insights {
		if insight.ID != idToRemove {
			filtered = append(filtered, insight)
		}
	}
	return filtered
}
