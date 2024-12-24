package state

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"zen/pkg/llm"
)

// NewPromptBuilder creates a new template builder instance
// It initializes empty sections and state data stores
func NewPromptBuilder(s *State) *PromptBuilder {
	return &PromptBuilder{
		state:     s,
		sections:  make([]PromptSection, 0),
		stateData: make(map[StateDataKey]interface{}),
		funcMap:   make(template.FuncMap),
	}
}

// Method to register template functions
func (tb *PromptBuilder) WithFunction(name string, fn interface{}) *PromptBuilder {
	if tb.err != nil {
		return tb
	}
	tb.funcMap[name] = fn
	return tb
}

// AddSection adds a new template section with the specified role
// Returns the builder for method chaining
func (tb *PromptBuilder) AddSection(role llm.Role, templateText string) *PromptBuilder {
	return tb.AddSectionWithName(role, templateText, "")
}

// AddSectionWithName adds a new template section with the specified role and name
// The name parameter allows identifying specific participants in the conversation
func (tb *PromptBuilder) AddSectionWithName(role llm.Role, templateText, name string) *PromptBuilder {
	if tb.err != nil {
		return tb
	}

	tb.sections = append(tb.sections, PromptSection{
		Role:     role,
		Template: templateText,
		Name:     name,
	})
	return tb
}

// Helper methods for common message types
// Each returns the builder for method chaining
func (tb *PromptBuilder) AddSystemSection(templateText string) *PromptBuilder {
	return tb.AddSection(llm.RoleSystem, templateText)
}

func (tb *PromptBuilder) AddUserSection(templateText string, name string) *PromptBuilder {
	return tb.AddSectionWithName(llm.RoleUser, templateText, name)
}

func (tb *PromptBuilder) AddAssistantSection(templateText string) *PromptBuilder {
	return tb.AddSection(llm.RoleAssistant, templateText)
}

func (tb *PromptBuilder) AddFunctionSection(templateText string, name string) *PromptBuilder {
	return tb.AddSectionWithName(llm.RoleFunction, templateText, name)
}

// WithManagerData adds a single piece of manager-provided data to the template context
// Returns an error if the specified key doesn't exist in the state's manager data
func (tb *PromptBuilder) WithManagerData(key StateDataKey) *PromptBuilder {
	if tb.err != nil {
		return tb
	}

	value, exists := tb.state.GetManagerData(key)
	if !exists {
		tb.err = fmt.Errorf("manager data for key %s not found", key)
		return tb
	}

	// Store the value with its original key
	tb.stateData[key] = value
	return tb
}

// WithManagerDataBatch adds multiple manager data keys at once
// Stops processing and returns error if any key is not found
func (tb *PromptBuilder) WithManagerDataBatch(keys ...StateDataKey) *PromptBuilder {
	for _, key := range keys {
		tb = tb.WithManagerData(key)
		if tb.err != nil {
			return tb
		}
	}
	return tb
}

// Compose processes all template sections and returns an array of formatted messages
// It combines state fields, manager data, and custom data for template rendering
func (tb *PromptBuilder) Compose() ([]llm.Message, error) {
	if tb.err != nil {
		return nil, tb.err
	}

	messages := make([]llm.Message, 0, len(tb.sections))

	for _, section := range tb.sections {
		// Create data map and automatically add all exported State fields
		data := make(map[string]interface{})

		// Add State fields
		stateValue := reflect.ValueOf(tb.state).Elem()
		stateType := stateValue.Type()
		for i := 0; i < stateValue.NumField(); i++ {
			field := stateType.Field(i)
			if field.IsExported() {
				data[field.Name] = stateValue.Field(i).Interface()
			}
		}

		// Add manager data with proper key names
		for k, v := range tb.stateData {
			data[string(k)] = v
		}

		// Add custom data
		for k, v := range tb.state.customData {
			data[k] = v
		}

		// Create and execute template
		tmpl, err := template.New("section").Funcs(tb.funcMap).Parse(section.Template)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template section: %w", err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("failed to execute template section (role=%s): %w", section.Role, err)
		}

		messages = append(messages, llm.Message{
			Role:    section.Role,
			Content: buf.String(),
			Name:    section.Name,
		})
	}

	return messages, nil
}
