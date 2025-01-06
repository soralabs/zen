package random_tools

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	toolkit "github.com/soralabs/toolkit/go"
)

type RandomNumberTool struct {
	toolkit.Tool

	history []RandomNumberGeneration
	mu      sync.Mutex
}

func NewRandomNumberTool() *RandomNumberTool {
	return &RandomNumberTool{}
}

func (t *RandomNumberTool) GetName() string {
	return "random_number"
}

func (t *RandomNumberTool) GetDescription() string {
	return "A tool that generates random numbers within a specified range"
}

func (t *RandomNumberTool) GetSchema() toolkit.Schema {
	return toolkit.Schema{
		Parameters: json.RawMessage(`{
            "type": "object",
            "required": ["min", "max"],
            "properties": {
                "min": {
                    "type": "number",
                    "description": "Minimum value (inclusive)"
                },
                "max": {
                    "type": "number",
                    "description": "Maximum value (inclusive)"
                }
            }
        }`),
		Returns: json.RawMessage(`{
            "type": "object",
            "properties": {
                "generated": {
                    "type": "number",
                    "description": "The generated random number"
                },
                "generation_id": {
                    "type": "integer",
                    "description": "ID of the generation in history"
                }
            }
        }`),
	}
}

func (t *RandomNumberTool) Execute(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Parse input parameters
	var input struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	}
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	if input.Min > input.Max {
		return nil, fmt.Errorf("min value cannot be greater than max value")
	}

	// Generate random number
	random := input.Min + rand.Float64()*(input.Max-input.Min)

	// Store result
	gen := RandomNumberGeneration{
		Min:       input.Min,
		Max:       input.Max,
		Generated: random,
		Time:      time.Now(),
	}
	t.history = append(t.history, gen)

	// Return result as JSON
	result := struct {
		Generated    float64 `json:"generated"`
		GenerationID int     `json:"generation_id"`
	}{
		Generated:    gen.Generated,
		GenerationID: len(t.history) - 1,
	}

	response, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response: %w", err)
	}

	return response, nil
}

func (t *RandomNumberTool) GetHistory() []RandomNumberGeneration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]RandomNumberGeneration{}, t.history...)
}

func (t *RandomNumberTool) GetGeneration(id int) (RandomNumberGeneration, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if id < 0 || id >= len(t.history) {
		return RandomNumberGeneration{}, fmt.Errorf("generation ID %d not found", id)
	}
	return t.history[id], nil
}
