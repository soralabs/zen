package engine

import (
	"fmt"
	"time"
	"zen/internal/db"
	"zen/internal/state"
	"zen/pkg/id"

	"github.com/pgvector/pgvector-go"
)

// NewState creates a new State instance with the provided options.
// It initializes an empty state and applies any provided configuration options.
func (e *Engine) NewStateFromFragment(fragment *db.Fragment) (*state.State, error) {
	state := state.NewState()
	state.Input = fragment
	if err := e.UpdateState(state); err != nil {
		return nil, fmt.Errorf("failed to create state: %w", err)
	}
	return state, nil
}

func (e *Engine) NewState(actorId, sessionId id.ID, input string) (*state.State, error) {
	embedding, err := e.llmClient.EmbedText(input)
	if err != nil {
		return nil, fmt.Errorf("failed to embed text: %w", err)
	}

	state := state.NewState()
	state.Input = &db.Fragment{
		ID:        id.New(),
		ActorID:   actorId,
		SessionID: sessionId,
		Content:   input,
		Metadata:  nil,
		Embedding: pgvector.NewVector(embedding),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := e.UpdateState(state); err != nil {
		return nil, fmt.Errorf("failed to create state: %w", err)
	}

	return state, nil
}

// UpdateState applies the provided options and collects manager data
func (e *Engine) UpdateState(s *state.State) error {
	// Reset manager data for fresh update
	// s.Reset()

	// NOTE THAT THE CURRENT MESSAGE IS NOT ADDED TO THE STATE, BUT AFTER MANAGERS HAVE PROVIDED THEIR CONTEXT
	// If we have managers configured, collect their context data
	if len(e.managers) > 0 {
		for _, manager := range e.managers {
			contextData, err := manager.Context(s)
			if err != nil {
				return fmt.Errorf("failed to get manager context: %w", err)
			}

			// log.Printf("Manager context: %v", contextData)

			// Add the manager's data to our state
			s.AddManagerData(contextData)
		}
	}

	// Add recent interactions to state
	recentInteractions, err := e.interactionFragmentStore.GetBySession(s.Input.SessionID, 20)
	if err != nil {
		return fmt.Errorf("failed to get recent interactions: %w", err)
	}

	s.RecentInteractions = append(s.RecentInteractions, recentInteractions...)

	actor, err := e.actorStore.GetByID(s.Input.ActorID)
	if err != nil {
		return fmt.Errorf("failed to get actor: %w", err)
	}
	s.Actor = actor

	return nil
}
