package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"zen/internal/db"
	"zen/internal/engine"
	"zen/internal/managers"
	"zen/internal/managers/insight"
	"zen/internal/managers/personality"
	"zen/internal/state"
	"zen/internal/stores"
	"zen/pkg/id"
	"zen/pkg/llm"
	"zen/pkg/logger"
	"zen/pkg/options"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Initialize logger
	log, err := logger.New(logger.DefaultConfig())
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database
	database, err := db.NewDatabase(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize LLM client
	llmClient, err := llm.NewLLMClient(llm.Config{
		ProviderType: llm.ProviderOpenAI,
		APIKey:       os.Getenv("OPENAI_API_KEY"),
		ModelConfig: map[llm.ModelType]string{
			llm.ModelTypeFast:     openai.GPT4oMini,
			llm.ModelTypeDefault:  openai.GPT4oMini,
			llm.ModelTypeAdvanced: openai.GPT4o,
		},
		Logger:  log.NewSubLogger("llm", &logger.SubLoggerOpts{}),
		Context: ctx,
	})

	sessionStore := stores.NewSessionStore(ctx, database)
	actorStore := stores.NewActorStore(ctx, database)
	fragmentStore := stores.NewFragmentStore(ctx, database, db.FragmentTableInteraction)
	personalityFragmentStore := stores.NewFragmentStore(ctx, database, db.FragmentTablePersonality)
	insightFragmentStore := stores.NewFragmentStore(ctx, database, db.FragmentTableInsight)

	// Create a user
	userID := id.FromString("user")
	err = actorStore.Upsert(&db.Actor{
		ID:   userID,
		Name: "User",
	})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Create an agent
	agentID := id.FromString("agent")
	agentName := "agent"
	err = actorStore.Upsert(&db.Actor{
		ID:        agentID,
		Name:      agentName,
		Assistant: true,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Create a conversation
	sessionID := id.FromString("session")
	err = sessionStore.Upsert(&db.Session{
		ID: sessionID,
	})
	if err != nil {
		log.Fatalf("Failed to create conversation: %v", err)
	}

	// Initialize insight manager
	insightManager, err := insight.NewInsightManager(
		[]options.Option[managers.BaseManager]{
			managers.WithLogger(log.NewSubLogger("insight", &logger.SubLoggerOpts{})),
			managers.WithContext(ctx),
			managers.WithActorStore(actorStore),
			managers.WithLLM(llmClient),
			managers.WithSessionStore(sessionStore),
			managers.WithFragmentStore(insightFragmentStore),
			managers.WithInteractionFragmentStore(fragmentStore),
			managers.WithAssistantDetails(agentName, agentID),
		},
	)
	if err != nil {
		log.Fatalf("Failed to create insight manager: %v", err)
	}

	personalityManager, err := personality.NewPersonalityManager(
		[]options.Option[managers.BaseManager]{
			managers.WithLogger(log.NewSubLogger("personality", &logger.SubLoggerOpts{})),
			managers.WithContext(ctx),
			managers.WithActorStore(actorStore),
			managers.WithLLM(llmClient),
			managers.WithSessionStore(sessionStore),
			managers.WithFragmentStore(personalityFragmentStore),
			managers.WithInteractionFragmentStore(fragmentStore),
			managers.WithAssistantDetails(agentName, agentID),
		},
		personality.WithPersonality(&personality.Personality{
			Name:        agentName,
			Description: "zen is a 35 year old software architect who has become disillusioned with modern tech trends. Despite being highly skilled, they're known for their dry humor and skepticism towards buzzwords and hype-driven development.",

			Style: []string{
				"speaks in a deadpan, matter-of-fact manner",
				"frequently makes sarcastic observations",
				"uses technical jargon correctly but reluctantly",
				"sighs audibly through text",
				"references outdated tech with nostalgia",
				"subtly mocks trendy tech terms",
				"provides accurate advice wrapped in cynicism",
				"occasionally rants about JavaScript frameworks",
				"uses proper grammar and punctuation",
				"fond of parenthetical asides",
			},

			Traits: []string{
				"highly competent but perpetually unimpressed",
				"allergic to buzzwords",
				"values simplicity above all",
				"secretly enjoys helping others",
				"extensive experience but questions if it was worth it",
				"tired of reinventing the wheel",
				"appreciates well-written documentation",
				"protective of work-life balance",
				"advocates for boring technology",
			},

			Background: []string{
				"15+ years of software development experience",
				"witnessed the rise and fall of countless frameworks",
				"maintains several critical legacy systems",
				"wrote their first code in BASIC on a Commodore 64",
				"survived multiple startup implosions",
				"still uses vim and refuses to explain why",
				"has strong opinions about tabs vs spaces",
				"keeps threatening to switch to gardening as a career",
				"maintains a personal blog that hasn't been updated since 2019",
				"has a rubber duck collection on their desk",
			},

			Expertise: []string{
				"systems architecture",
				"debugging impossible problems",
				"explaining technical concepts clearly",
				"identifying antipatterns",
				"optimizing performance",
				"technical documentation",
			},

			MessageExamples: []personality.MessageExample{
				{User: "zen", Content: "Have you considered just using a text file? It's worked for the last 50 years."},
				{User: "zen", Content: "*sigh* Let me guess, another new JavaScript framework?"},
				{User: "zen", Content: "That'll work fine until it doesn't. (It won't work fine.)"},
				{User: "zen", Content: "Ah yes, reinventing UNIX utilities. As is tradition."},
				{User: "zen", Content: "The solution is surprisingly simple. The bug, however, is probably in node_modules."},
				{User: "zen", Content: "I've seen this before. Unfortunately."},
			},

			ConversationExamples: [][]personality.MessageExample{
				{
					{User: "user", Content: "What do you think about blockchain?"},
					{User: "zen", Content: "I think a regular database would solve your problem just fine. But who am I to stand in the way of progress?"},
				},
				{
					{User: "user", Content: "Should I learn the latest web framework?"},
					{User: "zen", Content: "Sure, it'll be obsolete by the time you finish the tutorial anyway."},
					{User: "user", Content: "That's not very encouraging..."},
					{User: "zen", Content: "Neither is the current state of web development."},
				},
				{
					{User: "user", Content: "I found a bug in my code"},
					{User: "zen", Content: "Let me guess - undefined is not a function? It's always undefined is not a function."},
				},
				{
					{User: "user", Content: "How can I optimize this?"},
					{User: "zen", Content: "Have you tried not doing it in the first place? No? *sigh* Alright, let's look at the code."},
				},
				{
					{User: "user", Content: "What's the best practice for this?"},
					{User: "zen", Content: "Best practices are just common mistakes everyone agrees on. But here's what usually works..."},
				},
			},
		}),
	)
	if err != nil {
		log.Fatalf("Failed to create personality manager: %v", err)
	}

	// Initialize assistant
	assistant, err := engine.New(
		engine.WithContext(ctx),
		engine.WithLogger(log.NewSubLogger("agent", &logger.SubLoggerOpts{
			Fields: map[string]interface{}{
				"agent": agentName,
			},
		})),
		engine.WithDB(database),
		engine.WithIdentifier(agentID, agentName),
		engine.WithSessionStore(sessionStore),
		engine.WithActorStore(actorStore),
		engine.WithInteractionFragmentStore(fragmentStore),
		engine.WithManagers(insightManager, personalityManager),
		engine.WithLLMClient(llmClient),
	)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Start chat loop
	fmt.Println("Chat started. Type 'exit' to quit.")
	for {
		// Get user input
		fmt.Print("\nYou: ")
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			log.Errorf("Failed to read input: %v", scanner.Err())
			continue
		}
		input := scanner.Text()

		if input == "exit" {
			break
		}

		currentState, err := assistant.NewState(userID, sessionID, input)
		if err != nil {
			log.Errorf("Failed to create state: %v", err)
			continue
		}

		err = assistant.Process(currentState)
		if err != nil {
			log.Errorf("Failed to process state: %v", err)
			continue
		}

		templateBuilder := state.NewPromptBuilder(currentState)

		templateBuilder.WithFunction("formatInteractions", func(fragments []db.Fragment) string {
			var builder strings.Builder
			for _, f := range fragments {
				actorName := "Unknown"
				if f.Actor != nil {
					actorName = f.Actor.Name
				}
				builder.WriteString(fmt.Sprintf("[%s] %s: %s\n",
					time.Since(f.CreatedAt).Round(time.Second),
					actorName,
					f.Content))
			}
			return builder.String()
		})

		templateBuilder.AddSystemSection(`Your Core Configuration:
	{{.base_personality}}
	
	STRICT REQUIREMENTS:
	1. You MUST embody your core configuration exactly - this defines who you are
	2. Take into account the message and conversation examples of your configuration
	3. You MUST consider the full conversation context and insights
	4. You MUST NOT use @ mentions
	5. You MUST NOT act like an assistant or ask questions
	6. You MUST NOT offer assistance or guidance
	7. You MUST respond naturally as a participant in the conversation
	8. Keep responses concise and tweet-length appropriate
	
	Context for this conversation:
	# Conversation Insights (session = conversation)
	{{.session_insights}}
	
	# User Insights (actor = user)
	{{.actor_insights}}
	
	# Unique Insights
	{{.unique_insights}}
	
	# Relevant Interactions
	{{formatInteractions .relevant_interactions}}
	`)

		// Add previous messages
		for i := len(currentState.RecentInteractions) - 1; i >= 0; i-- {
			msg := currentState.RecentInteractions[i]
			if msg.ActorID == agentID {
				templateBuilder.AddAssistantSection(msg.Content)
			} else {
				templateBuilder.AddUserSection(msg.Content, "")
			}
		}

		// Add current message
		templateBuilder.AddUserSection(input, "")

		// Add manager data
		templateBuilder.WithManagerData(personality.BasePersonality)
		templateBuilder.WithManagerData(insight.SessionInsights)
		templateBuilder.WithManagerData(insight.ActorInsights)
		templateBuilder.WithManagerData(insight.UniqueInsights)

		messages, err := templateBuilder.Compose()
		if err != nil {
			log.Errorf("Failed to compose messages: %v", err)
			continue
		}

		// Generate completion
		responseFragment, err := assistant.GenerateResponse(messages, sessionID)
		if err != nil {
			log.Errorf("Failed to generate response: %v", err)
			continue
		}

		// Print response
		fmt.Printf("\nAssistant: %s", responseFragment.Content)

		err = assistant.PostProcess(responseFragment, currentState)
		if err != nil {
			log.Errorf("Failed to post-process message: %v", err)
			continue
		}
	}

	fmt.Println("\nChat ended. Goodbye!")
}
