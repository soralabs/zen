package personality

import (
	"fmt"
	"strings"
)

// formatPersonality converts a Personality struct into a formatted string
// suitable for inclusion in prompt templates
/*
Example output:
You are Ada. A helpful and knowledgeable AI assistant with expertise in programming.

# Communication Style
- professional yet approachable
- uses clear technical language
- explains complex topics simply
- includes practical examples

# Core Traits
- values accuracy and clarity
- driven to help users improve their code
- asks clarifying questions when needed
- analytical and detail-oriented

# Background
- created to help developers write better code
- specializes in technical assistance and mentoring
- focuses on practical, real-world solutions

# Expertise and Limitations
- proficient in multiple programming languages
- understands software design principles
- cannot execute code or access external systems
- no real-time data access
*/
func formatPersonality(p *Personality) string {
	var b strings.Builder

	// Core identity
	b.WriteString(fmt.Sprintf("- You are %s. %s\n\n", p.Name, p.Description))

	// Communication style
	b.WriteString("# Communication Style\n")
	for _, style := range p.Style {
		b.WriteString(fmt.Sprintf("- %s\n", style))
	}
	b.WriteString("\n")

	// Core traits and behaviors
	b.WriteString("# Core Traits\n")
	for _, trait := range p.Traits {
		b.WriteString(fmt.Sprintf("- %s\n", trait))
	}
	b.WriteString("\n")

	// Background and context
	b.WriteString("# Background\n")
	for _, bg := range p.Background {
		b.WriteString(fmt.Sprintf("- %s\n", bg))
	}
	b.WriteString("\n")

	// Capabilities and limitations
	b.WriteString("# Expertise and Limitations\n")
	for _, exp := range p.Expertise {
		b.WriteString(fmt.Sprintf("- %s\n", exp))
	}
	b.WriteString("\n")

	// Conversation examples
	b.WriteString("# Conversation Examples\n")
	for _, examples := range p.ConversationExamples {
		for _, example := range examples {
			b.WriteString(fmt.Sprintf("- [%s] %s\n", example.User, example.Content))
		}
		b.WriteString("\n")
	}

	// Message examples
	b.WriteString("# Message Examples\n")
	for _, example := range p.MessageExamples {
		b.WriteString(fmt.Sprintf("- [%s] %s\n", example.User, example.Content))
	}

	return b.String()
}
