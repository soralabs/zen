package personality

import (
	"fmt"
	"math/rand"
	"strings"
)

// randomSelect returns n random items from the input slice
// If n is greater than the length of the slice, returns all items in random order
func randomSelect(items []string, n int) []string {
	if n >= len(items) {
		// Create a copy to avoid modifying the original slice
		result := make([]string, len(items))
		copy(result, items)
		rand.Shuffle(len(result), func(i, j int) {
			result[i], result[j] = result[j], result[i]
		})
		return result
	}

	// Create an index slice
	indices := make([]int, len(items))
	for i := range indices {
		indices[i] = i
	}

	// Shuffle indices
	rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	// Select first n items using shuffled indices
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = items[indices[i]]
	}
	return result
}

// formatPersonality converts a Personality struct into a formatted string
// suitable for inclusion in prompt templates
/*
Example output:
You are Ada. A helpful and knowledgeable AI assistant with expertise in programming.

# Communication Style
- professional yet approachable
- uses clear technical language
- explains complex topics simply
- includes practical examples

# Core Traits
- values accuracy and clarity
- driven to help users improve their code
- asks clarifying questions when needed
- analytical and detail-oriented

# Background
- created to help developers write better code
- specializes in technical assistance and mentoring
- focuses on practical, real-world solutions

# Expertise and Limitations
- proficient in multiple programming languages
- understands software design principles
- cannot execute code or access external systems
- no real-time data access
*/
func formatPersonality(p *Personality) string {
	var b strings.Builder

	// Core identity
	b.WriteString(fmt.Sprintf("- You are %s. %s\n\n", p.Name, p.Description))

	// Communication style - randomly select 5
	b.WriteString("# Communication Style\n")
	for _, style := range randomSelect(p.Style, 5) {
		b.WriteString(fmt.Sprintf("- %s\n", style))
	}
	b.WriteString("\n")

	// Core traits and behaviors - randomly select 5
	b.WriteString("# Core Traits\n")
	for _, trait := range randomSelect(p.Traits, 5) {
		b.WriteString(fmt.Sprintf("- %s\n", trait))
	}
	b.WriteString("\n")

	// Background and context
	b.WriteString("# Background\n")
	for _, bg := range p.Background {
		b.WriteString(fmt.Sprintf("- %s\n", bg))
	}
	b.WriteString("\n")

	// Capabilities and limitations
	b.WriteString("# Expertise and Limitations\n")
	for _, exp := range p.Expertise {
		b.WriteString(fmt.Sprintf("- %s\n", exp))
	}
	b.WriteString("\n")

	// Conversation examples
	b.WriteString("# Conversation Examples\n")
	for _, examples := range p.ConversationExamples {
		for _, example := range examples {
			b.WriteString(fmt.Sprintf("- [%s] %s\n", example.User, example.Content))
		}
		b.WriteString("\n")
	}

	// Message examples
	b.WriteString("# Message Examples\n")
	for _, example := range p.MessageExamples {
		b.WriteString(fmt.Sprintf("- [%s] %s\n", example.User, example.Content))
	}

	return b.String()
}
