package engine

import (
	"context"
	"fmt"

	"github.com/soralabs/zen/id"
	"github.com/soralabs/zen/llm"
	"github.com/soralabs/zen/logger"
	"github.com/soralabs/zen/manager"
	"github.com/soralabs/zen/options"
	"github.com/soralabs/zen/stores"

	"gorm.io/gorm"
)

func (e *Engine) ValidateRequiredFields() error {
	if e.ctx == nil {
		return fmt.Errorf("context is required")
	}
	if e.db == nil {
		return fmt.Errorf("db is required")
	}
	if e.logger == nil {
		return fmt.Errorf("logger is required")
	}
	if e.actorStore == nil {
		return fmt.Errorf("actor store is required")
	}
	if e.sessionStore == nil {
		return fmt.Errorf("session store is required")
	}
	if e.interactionFragmentStore == nil {
		return fmt.Errorf("interaction fragment store is required")
	}
	if e.ID == "" {
		return fmt.Errorf("id is required")
	}
	if e.Name == "" {
		return fmt.Errorf("name is required")
	}
	if e.llmClient == nil {
		return fmt.Errorf("llm client is required")
	}
	return nil
}

func WithContext(ctx context.Context) options.Option[Engine] {
	return func(e *Engine) error {
		e.ctx = ctx
		return nil
	}
}

func WithDB(db *gorm.DB) options.Option[Engine] {
	return func(e *Engine) error {
		e.db = db
		return nil
	}
}

func WithLogger(logger *logger.Logger) options.Option[Engine] {
	return func(e *Engine) error {
		e.logger = logger
		return nil
	}
}

func WithIdentifier(id id.ID, name string) options.Option[Engine] {
	return func(e *Engine) error {
		e.ID = id
		e.Name = name
		return nil
	}
}

func WithInteractionFragmentStore(store *stores.FragmentStore) options.Option[Engine] {
	return func(e *Engine) error {
		e.interactionFragmentStore = store
		return nil
	}
}

func WithActorStore(store *stores.ActorStore) options.Option[Engine] {
	return func(e *Engine) error {
		e.actorStore = store
		return nil
	}
}

func WithSessionStore(store *stores.SessionStore) options.Option[Engine] {
	return func(e *Engine) error {
		e.sessionStore = store
		return nil
	}
}

func WithManagers(_managers ...manager.Manager) options.Option[Engine] {
	return func(e *Engine) error {
		// Create a map of available managers, checking for duplicates
		available := make(map[manager.ManagerID]manager.Manager)
		for _, m := range _managers {
			id := m.GetID()
			if _, exists := available[id]; exists {
				return fmt.Errorf("duplicate manager with ID %s", id)
			}
			available[id] = m
		}

		// Check dependencies
		for _, m := range _managers {
			for _, dep := range m.GetDependencies() {
				if _, ok := available[dep]; !ok {
					return fmt.Errorf("manager %s requires manager %s which was not provided",
						m.GetID(), dep)
				}
			}
		}

		e.managers = _managers
		return nil
	}
}

func WithManagerOrder(order []manager.ManagerID) options.Option[Engine] {
	return func(e *Engine) error {
		// Verify all managers in order exist
		managerMap := make(map[manager.ManagerID]bool)
		for _, m := range e.managers {
			managerMap[m.GetID()] = true
		}

		for _, id := range order {
			if !managerMap[id] {
				return fmt.Errorf("manager %s specified in order but not provided", id)
			}
		}

		e.managerOrder = order
		return nil
	}
}

func WithLLMClient(client *llm.LLMClient) options.Option[Engine] {
	return func(e *Engine) error {
		e.llmClient = client
		return nil
	}
}
