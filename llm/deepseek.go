package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/soralabs/zen/logger"
)

type DeepseekProvider struct {
	client *resty.Client
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
			ModelTypeFast:     "deepseek-chat",
			ModelTypeDefault:  "deepseek-chat",
			ModelTypeAdvanced: "deepseek-reasoner",
		}
	}

	// Role mapping
	roles := map[Role]string{
		RoleSystem:    "system",
		RoleUser:      "user",
		RoleAssistant: "assistant",
		RoleTool:      "tool",
	}

	client := resty.New().
		SetBaseURL("https://api.deepseek.com/v1").
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", config.DefaultProvider.APIKey)).
		SetHeader("Content-Type", "application/json").
		SetTimeout(2 * time.Minute)

	return &DeepseekProvider{
		client: client,
		models: models,
		logger: config.Logger,
		roles:  roles,
	}
}

type deepseekChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []deepseekMessage       `json:"messages"`
	Temperature      float32                 `json:"temperature,omitempty"`
	Tools            []deepseekTool          `json:"tools,omitempty"`
	ResponseFormat   *deepseekResponseFormat `json:"response_format,omitempty"`
	ReasoningContent string                  `json:"reasoning_content,omitempty"`
}

type deepseekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type deepseekTool struct {
	Type     string           `json:"type"`
	Function deepseekFunction `json:"function"`
}

type deepseekFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

type deepseekResponseFormat struct {
	Type string `json:"type"`
}

type deepseekChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role             string `json:"role"`
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// GenerateCompletion sends a conversation to the Deepseek Chat API
// and returns the model's text completion.
func (p *DeepseekProvider) GenerateCompletion(ctx context.Context, req CompletionRequest) (Message, error) {
	tools := make([]deepseekTool, len(req.Tools))
	for i, tool := range req.Tools {
		schema := tool.GetSchema()
		var params interface{}
		if err := json.Unmarshal(schema.Parameters, &params); err != nil {
			return Message{}, fmt.Errorf("failed to parse tool parameters: %w", err)
		}

		tools[i] = deepseekTool{
			Type: "function",
			Function: deepseekFunction{
				Name:        tool.GetName(),
				Description: tool.GetDescription(),
				Parameters:  params,
			},
		}
	}

	messages := p.convertMessages(req.Messages)
	chatReq := deepseekChatCompletionRequest{
		Model:       p.getModel(req.ModelType),
		Messages:    messages,
		Temperature: req.Temperature,
		ResponseFormat: &deepseekResponseFormat{
			Type: "text",
		},
	}
	if len(tools) > 0 {
		chatReq.Tools = tools
	}

	// Add reasoning_content for reasoner model
	if p.getModel(req.ModelType) == "deepseek-reasoner" {
		var reasoningContent string
		for _, msg := range req.Messages {
			if msg.Role == RoleUser {
				reasoningContent = msg.Content
				break
			}
		}
		chatReq.ReasoningContent = reasoningContent
	}

	var resp deepseekChatCompletionResponse
	httpResp, err := p.client.R().
		SetContext(ctx).
		SetBody(chatReq).
		SetResult(&resp).
		Post("/chat/completions")

	if err != nil {
		return Message{}, fmt.Errorf("Deepseek API error: %w", err)
	}

	if httpResp.StatusCode() != http.StatusOK {
		return Message{}, fmt.Errorf("Deepseek API error: %s", httpResp.String())
	}

	if len(resp.Choices) == 0 {
		return Message{}, fmt.Errorf("no completion returned")
	}

	choice := resp.Choices[0]
	if choice.Message.ReasoningContent != "" {
		p.logger.WithFields(map[string]interface{}{
			"reasoning": choice.Message.ReasoningContent,
		}).Info("Reasoning output")
	}
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

	chatReq := deepseekChatCompletionRequest{
		Model:       p.getModel(req.ModelType),
		Messages:    p.convertMessages(messages),
		Temperature: req.Temperature,
		ResponseFormat: &deepseekResponseFormat{
			Type: "json_object",
		},
	}

	var resp deepseekChatCompletionResponse
	httpResp, err := p.client.R().
		SetContext(ctx).
		SetBody(chatReq).
		SetResult(&resp).
		Post("/chat/completions")

	if err != nil {
		return fmt.Errorf("StructuredOutput Deepseek API error: %w", err)
	}

	if httpResp.StatusCode() != http.StatusOK {
		return fmt.Errorf("StructuredOutput Deepseek API error: %s", httpResp.String())
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("no completion returned")
	}

	return json.Unmarshal([]byte(resp.Choices[0].Message.Content), result)
}

// EmbedText generates an embedding vector for the given text using Deepseek's embedding model
func (p *DeepseekProvider) EmbedText(ctx context.Context, text string) ([]float32, error) {
	return nil, fmt.Errorf("embeddings not yet supported by Deepseek API")
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
func (p *DeepseekProvider) convertMessages(messages []Message) []deepseekMessage {
	converted := make([]deepseekMessage, len(messages))
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

		converted[i] = deepseekMessage{
			Role:    p.mapRole(msg.Role),
			Content: content,
			Name:    msg.Name,
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
