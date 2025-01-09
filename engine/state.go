package engine

import (
	"fmt"
	"time"

	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/id"
	"github.com/soralabs/zen/state"

	"github.com/pgvector/pgvector-go"
)

// NewState creates a new State instance with the provided options.
// It initializes an empty state and applies any provided configuration options.
func (e *Engine) PopulateStateFromFragment(state *state.State, fragment *db.Fragment) error {
	state.Input = fragment
	if err := e.UpdateState(state); err != nil {
		return fmt.Errorf("failed to create state: %w", err)
	}
	return nil
}

func (e *Engine) PopulateState(state *state.State, actorId, sessionId id.ID, input string) error {
	embedding, err := e.llmClient.EmbedText(input)
	if err != nil {
		return fmt.Errorf("failed to embed text: %w", err)
	}

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
		return fmt.Errorf("failed to create state: %w", err)
	}

	return nil
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

	relevantInteractions, err := e.interactionFragmentStore.SearchSimilar(s.Input.Embedding, s.Input.SessionID, 20)
	if err != nil {
		return fmt.Errorf("failed to get relevant interactions: %w", err)
	}

	s.RecentInteractions = append(s.RecentInteractions, recentInteractions...)

	s.RelevantInteractions = append(s.RelevantInteractions, e.filterInteractions(s.RecentInteractions, relevantInteractions)...)

	actor, err := e.actorStore.GetByID(s.Input.ActorID)
	if err != nil {
		return fmt.Errorf("failed to get actor: %w", err)
	}
	s.Actor = actor

	return nil
}

func (e *Engine) filterInteractions(i1, i2 []db.Fragment) []db.Fragment {
	seen := make(map[id.ID]bool, len(i2))

	for _, fragment := range i2 {
		seen[fragment.ID] = true
	}

	result := make([]db.Fragment, len(i2))
	copy(result, i2)

	for _, fragment := range i1 {
		if !seen[fragment.ID] {
			result = append(result, fragment)
			seen[fragment.ID] = true
		}
	}

	return result
}
