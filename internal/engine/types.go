package engine

import (
	"context"
	"zen/internal/managers"
	"zen/internal/stores"
	"zen/pkg/id"
	"zen/pkg/llm"
	"zen/pkg/logger"
	"zen/pkg/options"

	"gorm.io/gorm"
)

type Engine struct {
	options.RequiredFields

	ctx context.Context

	db *gorm.DB

	logger *logger.Logger

	ID   id.ID
	Name string

	// State management
	managers     []managers.Manager
	managerOrder []managers.ManagerID

	// stores
	actorStore   *stores.ActorStore
	sessionStore *stores.SessionStore

	interactionFragmentStore *stores.FragmentStore

	// LLM client
	llmClient *llm.LLMClient
}
