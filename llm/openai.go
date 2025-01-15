package llm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/soralabs/zen/logger"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type OpenAIProvider struct {
	client *openai.Client
	models map[ModelType]string
	logger *logger.Logger
	roles  map[Role]string
}

// NewOpenAIProvider creates and returns a new OpenAIProvider instance,
// initializing a default mapping for models and roles if none are provided.
func NewOpenAIProvider(config Config) *OpenAIProvider {
	// Default model mapping if not provided
	models := config.ModelConfig
	if models == nil {
		models = map[ModelType]string{
			ModelTypeFast:     openai.GPT4oMini,
			ModelTypeDefault:  openai.GPT4oMini,
			ModelTypeAdvanced: openai.GPT4o,
		}
	}

	// Role mapping
	roles := map[Role]string{
		RoleSystem:    openai.ChatMessageRoleSystem,
		RoleUser:      openai.ChatMessageRoleUser,
		RoleAssistant: openai.ChatMessageRoleAssistant,
		RoleTool:      openai.ChatMessageRoleTool,
	}

	return &OpenAIProvider{
		client: openai.NewClient(config.APIKey),
		models: models,
		logger: config.Logger,
		roles:  roles,
	}
}

// GenerateCompletion sends a conversation to the OpenAI ChatCompletion API
// and returns the model's text completion.
func (p *OpenAIProvider) GenerateCompletion(ctx context.Context, req CompletionRequest) (Message, error) {
	functions := make([]openai.FunctionDefinition, len(req.Tools))
	for i, tool := range req.Tools {
		schema := tool.GetSchema()
		functions[i] = openai.FunctionDefinition{
			Name:        tool.GetName(),
			Description: tool.GetDescription(),
			Parameters:  schema.Parameters,
		}
	}

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       p.getModel(req.ModelType),
		Messages:    p.convertMessages(req.Messages),
		Temperature: req.Temperature,
		Functions:   functions,
	})
	if err != nil {
		return Message{}, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return Message{}, fmt.Errorf("no completion returned")
	}

	// Handle function calls if present
	if resp.Choices[0].Message.FunctionCall != nil {
		call := resp.Choices[0].Message.FunctionCall
		for _, tool := range req.Tools {
			if tool.GetName() == call.Name {
				// Execute the tool
				result, err := tool.Execute(ctx, json.RawMessage(call.Arguments))
				if err != nil {
					return Message{}, fmt.Errorf("tool execution error: %w", err)
				}

				// Create a new message array with the tool result
				toolResultMessages := append(req.Messages,
					Message{
						Role:    RoleAssistant,
						Content: "",
						ToolCall: &ToolCall{
							Name:      call.Name,
							Arguments: string(call.Arguments),
						},
					},
					Message{
						Role:    RoleTool,
						Content: string(result),
						Name:    tool.GetName(),
					},
				)

				// Make a follow-up completion request with the tool result
				followUpReq := CompletionRequest{
					Messages:    toolResultMessages,
					ModelType:   req.ModelType,
					Temperature: req.Temperature,
					Tools:       req.Tools,
				}

				return p.GenerateCompletion(ctx, followUpReq)
			}
		}
		return Message{}, fmt.Errorf("function %s not found", call.Name)
	}

	return Message{
		Role:    RoleAssistant,
		Content: resp.Choices[0].Message.Content,
	}, nil
}

// GenerateStructuredOutput prompts the OpenAI API to return JSON data conforming
func (p *OpenAIProvider) GenerateStructuredOutput(ctx context.Context, req StructuredOutputRequest, result interface{}) error {
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		return fmt.Errorf("failed to generate schema: %w", err)
	}

	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    p.getModel(req.ModelType),
		Messages: p.convertMessages(req.Messages),
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   req.SchemaName,
				Schema: schema,
				Strict: req.StrictSchema,
			},
		},
		Temperature: req.Temperature,
	})
	if err != nil {
		return fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("no completion returned")
	}

	return schema.Unmarshal(resp.Choices[0].Message.Content, result)
}

// EmbedText generates an embedding vector for the given text using the Ada V2 model
func (p *OpenAIProvider) EmbedText(ctx context.Context, text string) ([]float32, error) {
	resp, err := p.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return resp.Data[0].Embedding, nil
}

// getModel returns the OpenAI model identifier for the given model type.
// Falls back to default model if type is not found.
func (p *OpenAIProvider) getModel(modelType ModelType) string {
	if model, ok := p.models[modelType]; ok {
		return model
	}
	return p.models[ModelTypeDefault]
}

// convertMessages transforms internal message format to OpenAI API format.
func (p *OpenAIProvider) convertMessages(messages []Message) []openai.ChatCompletionMessage {
	converted := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		converted[i] = openai.ChatCompletionMessage{
			Role:    p.mapRole(msg.Role),
			Content: msg.Content,
			Name:    msg.Name,
		}
		if msg.ToolCall != nil {
			converted[i].FunctionCall = &openai.FunctionCall{
				Name:      msg.ToolCall.Name,
				Arguments: msg.ToolCall.Arguments,
			}
		}
	}
	return converted
}

// mapRole converts internal role types to OpenAI API role strings.
func (p *OpenAIProvider) mapRole(role Role) string {
	if mappedRole, ok := p.roles[role]; ok {
		return mappedRole
	}
	return string(role)
}
