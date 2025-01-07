package engine

import (
	"context"

	"github.com/soralabs/zen/id"
	"github.com/soralabs/zen/llm"
	"github.com/soralabs/zen/logger"
	"github.com/soralabs/zen/manager"
	"github.com/soralabs/zen/options"
	"github.com/soralabs/zen/stores"

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
	managers     []manager.Manager
	managerOrder []manager.ManagerID

	// stores
	actorStore   *stores.ActorStore
	sessionStore *stores.SessionStore

	interactionFragmentStore *stores.FragmentStore

	// LLM client
	llmClient *llm.LLMClient
}
