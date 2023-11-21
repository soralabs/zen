package llm

import "context"

type Provider interface {
	GenerateCompletion(ctx context.Context, req CompletionRequest) (string, error)
	GenerateStructuredOutput(ctx context.Context, req StructuredOutputRequest, result interface{}) error
	EmbedText(ctx context.Context, text string) ([]float32, error)
}

type CompletionRequest struct {
	Messages    []Message
	ModelType   ModelType
	Temperature float32
}

type StructuredOutputRequest struct {
	Messages     []Message
	ModelType    ModelType
	Temperature  float32
	SchemaName   string
	StrictSchema bool
}
