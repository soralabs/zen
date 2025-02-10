package engine

import (
	"fmt"
	"time"

	"github.com/soralabs/zen/db"
	"github.com/soralabs/zen/id"
	"github.com/soralabs/zen/llm"
	"github.com/soralabs/zen/manager"

	"github.com/soralabs/zen/options"
	"github.com/soralabs/zen/state"

	"github.com/pgvector/pgvector-go"
	toolkit "github.com/soralabs/toolkit/go"
)

// New creates a new Engine instance with the provided options.
// Returns an error if required fields are missing or if actor creation fails.
func New(opts ...options.Option[Engine]) (*Engine, error) {
	e := &Engine{}
	if err := options.ApplyOptions(e, opts...); err != nil {
		return nil, fmt.Errorf("failed to create core: %w", err)
	}

	if err := e.UpsertActor(e.ID, e.Name, true); err != nil {
		return nil, fmt.Errorf("failed to upsert actor: %w", err)
	}

	return e, nil
}

// Process handles the processing of a new input through the runtime pipeline:
// 1. Retrieves actor and session information
// 2. Creates a copy of the input fragment
// 3. Executes all managers in parallel
// 4. Stores the processed input
// Returns an error if any step fails.
func (e *Engine) Process(currentState *state.State) error {
	var allMgrs []manager.ManagerID
	for _, m := range e.managers {
		allMgrs = append(allMgrs, m.GetID())
	}
	return e.ProcessWithFilter(currentState, allMgrs)
}

// ProcessWithFilter handles the processing of a new input through the runtime pipeline
// but only executes the managers whose IDs are in the provided managerFilter slice.
func (e *Engine) ProcessWithFilter(currentState *state.State, managerFilter []manager.ManagerID) error {
	return e.NewProcessBuilder().
		WithState(currentState).
		WithManagerFilter(managerFilter).
		Execute()
}

// PostProcess handles the post-processing of a response:
// 1. Retrieves actor and session information
// 2. Creates a copy of the response fragment
// 3. Executes all managers in sequence
// 4. Stores the processed response
// Returns an error if any step fails.
func (e *Engine) PostProcess(response *db.Fragment, currentState *state.State) error {
	var allMgrs []manager.ManagerID
	for _, m := range e.managers {
		allMgrs = append(allMgrs, m.GetID())
	}
	return e.PostProcessWithFilter(response, currentState, allMgrs)
}

// PostProcessWithFilter handles the post-processing of a response but only executes the managers specified in managerFilter.
func (e *Engine) PostProcessWithFilter(response *db.Fragment, currentState *state.State, managerFilter []manager.ManagerID) error {
	return e.NewPostProcessBuilder().
		WithState(currentState).
		WithResponse(response).
		WithManagerFilter(managerFilter).
		Execute()
}

// GenerateResponse creates a new response using the LLM:
// 1. Generates completion from provided messages
// 2. Creates embedding for the response
// 3. Builds response fragment with metadata
// Returns the response fragment and any error encountered.
func (e *Engine) GenerateResponse(messages []llm.Message, sessionID id.ID, tools ...toolkit.Tool) (*db.Fragment, error) {
	// Generate completion
	response, err := e.llmClient.GenerateCompletion(llm.CompletionRequest{
		Messages:    messages,
		ModelType:   llm.ModelTypeDefault,
		Temperature: 0.7,
		Tools:       tools,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate completion: %v", err)
	}

	// Generate embedding for the response
	embedding, err := e.llmClient.EmbedText(response.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedding for response: %v", err)
	}

	// Create response fragment
	responseFragment := &db.Fragment{
		ID:        id.New(),
		ActorID:   e.ID,
		SessionID: sessionID,
		Content:   response.Content,
		Embedding: pgvector.NewVector(embedding),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  nil,
	}

	return responseFragment, nil
}

// StartBackgroundProcesses initiates background processes for all managers.
// Each manager's background process runs in its own goroutine.
func (e *Engine) StartBackgroundProcesses() {
	for _, manager := range e.managers {
		go manager.StartBackgroundProcesses()
	}
}

// StopBackgroundProcesses terminates background processes for all managers.
func (e *Engine) StopBackgroundProcesses() {
	for _, manager := range e.managers {
		manager.StopBackgroundProcesses()
	}
}

// AddManager adds a new manager to the runtime.
// Validates that:
// 1. The manager ID is not duplicate
// 2. All manager dependencies are available
// Returns an error if validation fails.
func (e *Engine) AddManager(newManager manager.Manager) error {
	// Check for duplicates
	for _, m := range e.managers {
		if m.GetID() == newManager.GetID() {
			return fmt.Errorf("duplicate manager with ID %s", newManager.GetID())
		}
	}

	// Check dependencies
	available := make(map[manager.ManagerID]bool)
	for _, m := range e.managers {
		available[m.GetID()] = true
	}

	for _, dep := range newManager.GetDependencies() {
		if _, exists := available[dep]; !exists {
			return fmt.Errorf("manager %s requires manager %s which was not provided",
				newManager.GetID(), dep)
		}
	}

	e.managers = append(e.managers, newManager)
	return nil
}

// executeManagersInOrder runs managers in a specified order:
// 1. Creates a map for quick manager lookup
// 2. Uses managerOrder if specified, otherwise uses registration order
// 3. Executes each manager with the provided function
// Returns an error if any manager execution fails.
func (e *Engine) executeManagersInOrder(currentState *state.State, executeFn func(manager.Manager) error) error {
	// Create a map for quick manager lookup
	managerMap := make(map[manager.ManagerID]manager.Manager)
	for _, m := range e.managers {
		managerMap[m.GetID()] = m
	}

	// If no order specified, use registration order
	executionOrder := e.managerOrder
	if len(executionOrder) == 0 {
		executionOrder = make([]manager.ManagerID, len(e.managers))
		for i, m := range e.managers {
			executionOrder[i] = m.GetID()
		}
	}

	// Execute managers in specified order
	for _, managerID := range executionOrder {
		if manager, exists := managerMap[managerID]; exists {
			if err := executeFn(manager); err != nil {
				return fmt.Errorf("manager %s failed: %w", managerID, err)
			}
		}
	}

	return nil
}
