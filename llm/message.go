package llm

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type ToolCall struct {
	Name      string
	Arguments string
}

type Message struct {
	Role     Role
	Content  string
	Name     string
	ToolCall *ToolCall
}

type ModelType string

const (
	ModelTypeFast     ModelType = "fast"
	ModelTypeDefault  ModelType = "default"
	ModelTypeAdvanced ModelType = "advanced"
)
