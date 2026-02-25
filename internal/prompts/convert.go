// =============================================================================
// FILE: internal/prompts/convert.go
// PURPOSE: Prompt conversion helpers. Converts between prompt selection values
//          and internal types (actions, areas, etc.).
//          Ports Python prompts/promptConvert.py.
// =============================================================================

package prompts

import (
	"strings"
)

// ---------------------------------------------------------------------------
// Conversion helpers
// ---------------------------------------------------------------------------

// ActionFromLabel converts a display label to an action name.
//
// Parameters:
//   - label: The display label (e.g. "Download").
//
// Returns:
//   - The lowercase action name (e.g. "download").
func ActionFromLabel(label string) string {
	return strings.ToLower(strings.TrimSpace(label))
}

// AreasFromLabels converts display labels to area names.
//
// Parameters:
//   - labels: The display labels.
//
// Returns:
//   - Slice of lowercase area names.
func AreasFromLabels(labels []string) []string {
	var areas []string
	for _, l := range labels {
		area := strings.ToLower(strings.TrimSpace(l))
		if area == "all" {
			return AllAreas()
		}
		if area != "" {
			areas = append(areas, area)
		}
	}
	return areas
}

// LabelFromAction converts an action name to a display label.
//
// Parameters:
//   - action: The action name.
//
// Returns:
//   - The display label with first letter capitalized.
func LabelFromAction(action string) string {
	if action == "" {
		return ""
	}
	return strings.ToUpper(action[:1]) + action[1:]
}

// LabelsFromAreas converts area names to display labels.
//
// Parameters:
//   - areas: The area names.
//
// Returns:
//   - Slice of display labels.
func LabelsFromAreas(areas []string) []string {
	var labels []string
	for _, a := range areas {
		labels = append(labels, LabelFromAction(a))
	}
	return labels
}
