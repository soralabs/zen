package engine

import (
	"context"

	"github.com/soralabs/zen/pkg/id"
	"github.com/soralabs/zen/pkg/llm"
	"github.com/soralabs/zen/pkg/logger"
	"github.com/soralabs/zen/pkg/manager"
	"github.com/soralabs/zen/pkg/options"
	"github.com/soralabs/zen/pkg/stores"

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
