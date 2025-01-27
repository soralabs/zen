package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
	"github.com/soralabs/zen/logger"
)

type DeepseekProvider struct {
	client *deepseek.Client
	models map[ModelType]string
	logger *logger.Logger
	roles  map[Role]string
}

// NewDeepseekProvider creates and returns a new DeepseekProvider instance,
// initializing a default mapping for models and roles if none are provided.
func NewDeepseekProvider(config Config) *DeepseekProvider {
	// Default model mapping if not provided
	models := config.DefaultProvider.ModelConfig
	if models == nil {
		models = map[ModelType]string{
			ModelTypeFast:     deepseek.DeepSeekChat,
			ModelTypeDefault:  deepseek.DeepSeekChat,
			ModelTypeAdvanced: deepseek.DeepSeekReasoner,
		}
	}

	// Role mapping
	roles := map[Role]string{
		RoleSystem:    constants.ChatMessageRoleSystem,
		RoleUser:      constants.ChatMessageRoleUser,
		RoleAssistant: constants.ChatMessageRoleAssistant,
		RoleTool:      "tool",
	}

	return &DeepseekProvider{
		client: deepseek.NewClient(config.DefaultProvider.APIKey),
		models: models,
		logger: config.Logger,
		roles:  roles,
	}
}

// GenerateCompletion sends a conversation to the Deepseek Chat API
// and returns the model's text completion.
func (p *DeepseekProvider) GenerateCompletion(ctx context.Context, req CompletionRequest) (Message, error) {
	tools := make([]deepseek.Tools, len(req.Tools))
	for i, tool := range req.Tools {
		schema := tool.GetSchema()
		params := &deepseek.Parameters{}
		if err := json.Unmarshal(schema.Parameters, params); err != nil {
			return Message{}, fmt.Errorf("failed to parse tool parameters: %w", err)
		}

		tools[i] = deepseek.Tools{
			Type: "function",
			Function: deepseek.Function{
				Name:        tool.GetName(),
				Description: tool.GetDescription(),
				Parameters:  params,
			},
		}
	}

	messages := p.convertMessages(req.Messages)
	chatReq := deepseek.ChatCompletionRequest{
		Model:       p.getModel(req.ModelType),
		Messages:    messages,
		Temperature: req.Temperature,
	}
	if len(tools) > 0 {
		chatReq.Tools = tools
	}

	resp, err := p.client.CreateChatCompletion(ctx, &chatReq)
	if err != nil {
		return Message{}, fmt.Errorf("Deepseek API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return Message{}, fmt.Errorf("no completion returned")
	}

	choice := resp.Choices[0]
	toolCall := p.extractToolCall(choice.Message.Content)
	if toolCall != nil {
		for _, tool := range req.Tools {
			if tool.GetName() == toolCall.Name {
				// Execute the tool
				result, err := tool.Execute(ctx, json.RawMessage(toolCall.Arguments))
				if err != nil {
					return Message{}, fmt.Errorf("tool execution error: %w", err)
				}

				// Create a new message array with the tool result
				toolResultMessages := append(req.Messages,
					Message{
						Role:    RoleAssistant,
						Content: "",
						ToolCall: &ToolCall{
							Name:      toolCall.Name,
							Arguments: toolCall.Arguments,
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
		return Message{}, fmt.Errorf("function %s not found", toolCall.Name)
	}

	return Message{
		Role:    RoleAssistant,
		Content: choice.Message.Content,
	}, nil
}

// extractToolCall extracts tool call information from a message
func (p *DeepseekProvider) extractToolCall(content string) *ToolCall {
	// Parse the content as JSON to check for function call
	var parsed struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	}
	if err := json.Unmarshal([]byte(content), &parsed); err == nil && parsed.Name != "" {
		return &ToolCall{
			Name:      parsed.Name,
			Arguments: parsed.Arguments,
		}
	}
	return nil
}

// GenerateStructuredOutput prompts the Deepseek API to return JSON data
func (p *DeepseekProvider) GenerateStructuredOutput(ctx context.Context, req StructuredOutputRequest, result interface{}) error {
	resultType := reflect.TypeOf(result)
	if resultType.Kind() == reflect.Ptr {
		resultType = resultType.Elem()
	}
	description := fmt.Sprintf("The response should be a valid JSON object matching the following Go struct type: %s", resultType.String())

	systemMessage := Message{
		Role:    RoleSystem,
		Content: description,
	}
	messages := append([]Message{systemMessage}, req.Messages...)

	chatReq := deepseek.ChatCompletionRequest{
		Model:       p.getModel(req.ModelType),
		Messages:    p.convertMessages(messages),
		Temperature: req.Temperature,
		ResponseFormat: &deepseek.ResponseFormat{
			Type: "json_object",
		},
	}

	resp, err := p.client.CreateChatCompletion(ctx, &chatReq)
	if err != nil {
		return fmt.Errorf("StructuredOutput Deepseek API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("no completion returned")
	}

	return json.Unmarshal([]byte(resp.Choices[0].Message.Content), result)
}

// EmbedText generates an embedding vector for the given text using Deepseek's embedding model
func (p *DeepseekProvider) EmbedText(ctx context.Context, text string) ([]float32, error) {
	// TODO: Implement once embedding support is added to deepseek-go
	return nil, fmt.Errorf("embeddings not yet supported by deepseek-go")
}

// getModel returns the Deepseek model identifier for the given model type.
// Falls back to default model if type is not found.
func (p *DeepseekProvider) getModel(modelType ModelType) string {
	if model, ok := p.models[modelType]; ok {
		return model
	}
	return p.models[ModelTypeDefault]
}

// convertMessages transforms internal message format to Deepseek API format.
func (p *DeepseekProvider) convertMessages(messages []Message) []deepseek.ChatCompletionMessage {
	converted := make([]deepseek.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		content := msg.Content
		if msg.ToolCall != nil {
			// Format tool call as JSON in the content
			toolCall := map[string]interface{}{
				"name":      msg.ToolCall.Name,
				"arguments": msg.ToolCall.Arguments,
			}
			if jsonContent, err := json.Marshal(toolCall); err == nil {
				content = string(jsonContent)
			}
		}

		converted[i] = deepseek.ChatCompletionMessage{
			Role:    p.mapRole(msg.Role),
			Content: content,
		}
	}
	return converted
}

// mapRole converts internal role types to Deepseek API role strings.
func (p *DeepseekProvider) mapRole(role Role) string {
	if mappedRole, ok := p.roles[role]; ok {
		return mappedRole
	}
	return string(role)
}
