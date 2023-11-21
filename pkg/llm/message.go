package llm

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleFunction  Role = "function"
)

type Message struct {
	Role    Role
	Content string
	Name    string
}

type ModelType string

const (
	ModelTypeFast     ModelType = "fast"
	ModelTypeDefault  ModelType = "default"
	ModelTypeAdvanced ModelType = "advanced"
)
