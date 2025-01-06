package random_tools

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	toolkit "github.com/soralabs/toolkit/go"
)

const (
	alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	alphabetic   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numeric      = "0123456789"
)

type RandomStringTool struct {
	toolkit.Tool

	history []RandomStringGeneration
	mu      sync.Mutex
}

func NewRandomStringTool() *RandomStringTool {
	return &RandomStringTool{}
}

func (t *RandomStringTool) GetSchema() toolkit.Schema {
	return toolkit.Schema{
		Parameters: json.RawMessage(`{
			"type": "object",
			"required": ["length"],
			"properties": {
				"length": {
					"type": "integer",
					"description": "Length of the string to generate",
					"minimum": 1,
					"maximum": 1000
				},
				"charset": {
					"type": "string",
					"description": "Character set to use (alphanumeric, alphabetic, numeric, custom)",
					"default": "alphanumeric",
					"enum": ["alphanumeric", "alphabetic", "numeric", "custom"]
				},
				"custom_charset": {
					"type": "string",
					"description": "Custom character set to use when charset is set to 'custom'",
					"minLength": 1
				}
			}
		}`),
		Returns: json.RawMessage(`{
			"type": "string",
			"description": "Generated random string"
		}`),
	}
}

func (t *RandomStringTool) Execute(params json.RawMessage) (interface{}, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	var input struct {
		Length        int    `json:"length"`
		Charset       string `json:"charset"`
		CustomCharset string `json:"custom_charset,omitempty"`
	}
	if err := json.Unmarshal(params, &input); err != nil {
		return nil, fmt.Errorf("failed to parse parameters: %w", err)
	}

	if input.Length <= 0 {
		return nil, fmt.Errorf("length must be greater than 0")
	}

	// Determine which charset to use
	var chars string
	switch input.Charset {
	case "alphanumeric", "":
		chars = alphanumeric
	case "alphabetic":
		chars = alphabetic
	case "numeric":
		chars = numeric
	case "custom":
		if input.CustomCharset == "" {
			return nil, fmt.Errorf("custom charset cannot be empty")
		}
		chars = input.CustomCharset
	default:
		return nil, fmt.Errorf("invalid charset specified")
	}

	// Generate random string
	result := make([]byte, input.Length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}

	// Store in history
	t.history = append(t.history, RandomStringGeneration{
		Length:  input.Length,
		Charset: input.Charset,
		Result:  string(result),
		Time:    time.Now(),
	})

	// Trim history if it gets too long
	if len(t.history) > 100 {
		t.history = t.history[1:]
	}

	return string(result), nil
}

func (t *RandomStringTool) GetHistory() interface{} {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.history
}

func (t *RandomStringTool) ClearHistory() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.history = nil
}
