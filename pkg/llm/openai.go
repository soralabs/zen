package llm

import (
	"context"
	"fmt"
	"zen/pkg/logger"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type OpenAIProvider struct {
	client *openai.Client
	models map[ModelType]string
	logger *logger.Logger
	roles  map[Role]string
}

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
		RoleFunction:  openai.ChatMessageRoleFunction,
	}

	return &OpenAIProvider{
		client: openai.NewClient(config.APIKey),
		models: models,
		logger: config.Logger,
		roles:  roles,
	}
}

func (p *OpenAIProvider) GenerateCompletion(ctx context.Context, req CompletionRequest) (string, error) {
	resp, err := p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       p.getModel(req.ModelType),
		Messages:    p.convertMessages(req.Messages),
		Temperature: req.Temperature,
	})
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no completion returned")
	}

	return resp.Choices[0].Message.Content, nil
}

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
