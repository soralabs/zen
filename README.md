# Zen - An Advanced Conversational LLM Framework

## Table of Contents
- [Overview](#overview)
- [Core Features](#core-features)
- [Extension Points](#extension-points)
- [Quick Start](#quick-start)

## Overview
Zen is a highly modular conversation engine built in Go that emphasizes pluggable architecture and platform independence. It provides a flexible foundation for building conversational systems through:

- Plugin-based architecture with hot-swappable components
- Multi-provider LLM support (OpenAI, custom providers)
- Cross-platform conversation management
- Extensible manager system for custom behaviors
- Vector-based semantic storage with pgvector

## Core Features

### Plugin Architecture
- **Manager System**: Extend functionality through custom managers
  - Insight Manager: Extracts and maintains conversation insights
  - Personality Manager: Handles response behavior and style
  - Custom Managers: Add your own specialized behaviors

### State Management
- **Shared State System**: Centralized state management across components
  - Manager-specific data storage
  - Custom data injection
  - Cross-manager communication

### LLM Integration
- **Provider Abstraction**: Support for multiple LLM providers
  - Built-in OpenAI support
  - Extensible provider interface for custom LLMs
  - Configurable model selection per operation
  - Automatic fallback and retry handling

### Platform Support
- **Platform Agnostic Core**: 
  - Abstract conversation engine independent of platforms
  - Built-in support for CLI chat and Twitter
  - Extensible platform manager interface
  - Example implementations for new platform integration

### Storage Layer
- **Flexible Data Storage**:
  - PostgreSQL with pgvector for semantic search
  - GORM-based data models
  - Customizable fragment storage
  - Vector embedding support

## Extension Points
1. **LLM Providers**: Add new AI providers by implementing the LLM interface
```go
type Provider interface {
    GenerateCompletion(context.Context, CompletionRequest) (string, error)
    GenerateJSON(context.Context, JSONRequest, interface{}) error
    EmbedText(context.Context, string) ([]float32, error)
}
```

2. **Managers**: Create new behaviors by implementing the Manager interface
```go
type Manager interface {
    GetID() ManagerID
    GetDependencies() []ManagerID
    Process(state *state.State) error
    PostProcess(state *state.State) error
    Context(state *state.State) ([]state.StateData, error)
    Store(fragment *db.Fragment) error
    StartBackgroundProcesses()
    StopBackgroundProcesses()
    RegisterEventHandler(callback EventCallbackFunc)
    triggerEvent(eventData EventData)
}
```

## Quick Start
1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Install dependencies: `go mod download`
4. Run the chat example: `go run examples/chat/main.go`
5. Run the Twitter bot: `go run examples/twitter/main.go`

## Environment Variables
```env
DB_URL=postgresql://user:password@localhost:5432/zen
OPENAI_API_KEY=your_openai_api_key

Platform-specific credentials as needed
```

## Architecture
The project follows a clean, modular architecture:

- `internal/engine`: Core conversation engine
- `internal/managers`: Plugin manager system
- `internal/state`: Shared state management
- `pkg/llm`: LLM provider interfaces
- `examples/`: Reference implementations
