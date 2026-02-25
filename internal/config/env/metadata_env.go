// =============================================================================
// FILE: internal/config/env/metadata_env.go
// PURPOSE: Metadata action environment variable defaults.
//          Ports Python of_env/values/action/metadata.py.
// =============================================================================

package env

// MetadataEnabled returns whether metadata operations are enabled.
func MetadataEnabled() bool {
	return GetBool("OF_METADATA_ENABLED", true)
}
