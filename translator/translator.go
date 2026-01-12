package translator

import "strings"

// Translator defines the interface for converting flags between tools
type Translator interface {
	// Name returns the translator identifier (e.g., "ls2eza")
	Name() string

	// SourceTool returns the name of the source tool (e.g., "ls")
	SourceTool() string

	// TargetTool returns the name of the target tool (e.g., "eza")
	TargetTool() string

	// Translate converts source tool arguments to target tool arguments
	Translate(args []string) []string
}

// EnvVarName returns the environment variable name for mode override
// Format: REFLAG_<NAME>_MODE (e.g., REFLAG_LS2EZA_MODE)
func EnvVarName(t Translator) string {
	return "REFLAG_" + strings.ToUpper(t.Name()) + "_MODE"
}
