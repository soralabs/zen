package engine

import (
	"fmt"

	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/id"
	"github.com/soralabs/zen/manager"
)

func (e *Engine) UpsertSession(sessionID id.ID) error {
	err := e.sessionStore.Upsert(&db.Session{
		ID: sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (e *Engine) UpsertActor(actorID id.ID, actorName string, assistant bool) error {
	err := e.actorStore.Upsert(&db.Actor{
		ID:        actorID,
		Name:      actorName,
		Assistant: assistant,
	})
	if err != nil {
		return fmt.Errorf("failed to create actor: %w", err)
	}

	return nil
}

func (e *Engine) UpsertInteractionFragment(fragment *db.Fragment) error {
	return e.interactionFragmentStore.Upsert(fragment)
}

func (e *Engine) DoesInteractionFragmentExist(fragmentID id.ID) (bool, error) {
	fragment, err := e.interactionFragmentStore.GetByID(fragmentID)
	if err != nil {
		return false, err
	}
	return fragment != nil, nil
}

// Helper function to check if a manager ID is in a list
func (e *Engine) containsManager(filter []manager.ManagerID, id manager.ManagerID) bool {
	for _, mid := range filter {
		if mid == id {
			return true
		}
	}
	return false
}
