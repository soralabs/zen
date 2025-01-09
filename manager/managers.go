package manager

import "github.com/soralabs/zen/id"

// Standard manager identifiers used throughout the system
const (
	// BaseManagerID represents the core manager implementation
	BaseManagerID id.ManagerID = "base"

	// PersonalityManagerID handles agent personality and behavior traits
	PersonalityManagerID id.ManagerID = "personality"

	// InsightManagerID manages conversation and user insights
	InsightManagerID id.ManagerID = "insight"

	// TwitterManagerID handles Twitter-specific interactions
	TwitterManagerID id.ManagerID = "twitter"
)
