package twitter

import (
	"fmt"
	"time"

	"github.com/soralabs/zen/pkg/db"
	"github.com/soralabs/zen/pkg/engine"
	"github.com/soralabs/zen/pkg/id"
	"github.com/soralabs/zen/pkg/logger"
	"github.com/soralabs/zen/pkg/manager"
	"github.com/soralabs/zen/pkg/managers/insight"
	"github.com/soralabs/zen/pkg/managers/personality"
	twitter_manager "github.com/soralabs/zen/pkg/managers/twitter"
	"github.com/soralabs/zen/pkg/options"
	"github.com/soralabs/zen/pkg/stores"
	"github.com/soralabs/zen/pkg/twitter"
)

func New(opts ...options.Option[Twitter]) (*Twitter, error) {
	k := &Twitter{
		stopChan: make(chan struct{}),
		twitterConfig: TwitterConfig{
			MonitorInterval: IntervalConfig{
				Min: 60 * time.Second,
				Max: 120 * time.Second,
			}, // default interval
		},
	}

	// Apply options
	if err := options.ApplyOptions(k, opts...); err != nil {
		return nil, fmt.Errorf("failed to apply options: %w", err)
	}

	// Validate required fields
	if err := k.ValidateRequiredFields(); err != nil {
		return nil, err
	}

	// Initialize Twitter client if enabled
	if k.twitterConfig.Credentials.CT0 == "" || k.twitterConfig.Credentials.AuthToken == "" {
		return nil, fmt.Errorf("Twitter credentials required when Twitter is enabled")
	}

	k.twitterClient = twitter.NewClient(
		k.ctx,
		k.logger.NewSubLogger("twitter", &logger.SubLoggerOpts{}),
		twitter.TwitterCredential{
			CT0:       k.twitterConfig.Credentials.CT0,
			AuthToken: k.twitterConfig.Credentials.AuthToken,
		},
	)

	// Create agent
	if err := k.create(); err != nil {
		return nil, err
	}

	return k, nil
}

func (k *Twitter) Start() error {
	go k.monitorTwitter()
	return nil
}

func (k *Twitter) Stop() error {
	return nil
}

func (k *Twitter) create() error {
	// Initialize stores
	sessionStore := stores.NewSessionStore(k.ctx, k.database)
	actorStore := stores.NewActorStore(k.ctx, k.database)
	interactionFragmentStore := stores.NewFragmentStore(k.ctx, k.database, db.FragmentTableInteraction)
	personalityFragmentStore := stores.NewFragmentStore(k.ctx, k.database, db.FragmentTablePersonality)
	insightFragmentStore := stores.NewFragmentStore(k.ctx, k.database, db.FragmentTableInsight)
	twitterFragmentStore := stores.NewFragmentStore(k.ctx, k.database, db.FragmentTableTwitter)

	assistantName := "zen"
	assistantID := id.FromString("zen")

	// Initialize insight manager
	insightManager, err := insight.NewInsightManager(
		[]options.Option[manager.BaseManager]{
			manager.WithLogger(k.logger.NewSubLogger("insight", &logger.SubLoggerOpts{})),
			manager.WithContext(k.ctx),
			manager.WithActorStore(actorStore),
			manager.WithLLM(k.llmClient),
			manager.WithSessionStore(sessionStore),
			manager.WithFragmentStore(insightFragmentStore),
			manager.WithInteractionFragmentStore(interactionFragmentStore),
			manager.WithAssistantDetails(assistantName, assistantID),
		},
	)
	if err != nil {
		return err
	}

	personalityManager, err := personality.NewPersonalityManager(
		[]options.Option[manager.BaseManager]{
			manager.WithLogger(k.logger.NewSubLogger("personality", &logger.SubLoggerOpts{})),
			manager.WithContext(k.ctx),
			manager.WithActorStore(actorStore),
			manager.WithLLM(k.llmClient),
			manager.WithSessionStore(sessionStore),
			manager.WithFragmentStore(personalityFragmentStore),
			manager.WithInteractionFragmentStore(interactionFragmentStore),
			manager.WithAssistantDetails(assistantName, assistantID),
		},
		personality.WithPersonality(&personality.Personality{
			Name:        assistantName,
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
		return err
	}

	// Initialize assistant
	assistant, err := engine.New(
		engine.WithContext(k.ctx),
		engine.WithLogger(k.logger.NewSubLogger("agent", &logger.SubLoggerOpts{
			Fields: map[string]interface{}{
				"agent": "zen",
			},
		})),
		engine.WithDB(k.database),
		engine.WithIdentifier(assistantID, assistantName),
		engine.WithSessionStore(sessionStore),
		engine.WithActorStore(actorStore),
		engine.WithInteractionFragmentStore(interactionFragmentStore),
		engine.WithManagers(insightManager, personalityManager),
	)
	if err != nil {
		return err
	}

	twitterManager, err := twitter_manager.NewTwitterManager(
		[]options.Option[manager.BaseManager]{
			manager.WithLogger(k.logger.NewSubLogger("twitter", &logger.SubLoggerOpts{})),
			manager.WithContext(k.ctx),
			manager.WithActorStore(actorStore),
			manager.WithLLM(k.llmClient),
			manager.WithSessionStore(sessionStore),
			manager.WithFragmentStore(twitterFragmentStore),
			manager.WithInteractionFragmentStore(interactionFragmentStore),
			manager.WithAssistantDetails(assistantName, assistantID),
		},
		twitter_manager.WithTwitterClient(
			k.twitterClient,
		),
		twitter_manager.WithTwitterUsername(
			k.twitterConfig.Credentials.User,
		),
	)
	if err != nil {
		return err
	}

	if err := assistant.AddManager(twitterManager); err != nil {
		return err
	}

	k.assistant = assistant

	return nil
}
