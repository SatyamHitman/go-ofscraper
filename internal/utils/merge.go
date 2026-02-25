// =============================================================================
// FILE: internal/utils/merge.go
// PURPOSE: Merge utility functions. Provides helpers for merging slices, maps,
//          and data structures used during post/media aggregation from multiple
//          API sources. Ports Python utils/merge.py.
// =============================================================================

package utils

// ---------------------------------------------------------------------------
// Slice merge/dedup
// ---------------------------------------------------------------------------

// MergeStringSlices merges multiple string slices into one, removing duplicates.
// Order of first occurrence is preserved.
//
// Parameters:
//   - slices: The string slices to merge.
//
// Returns:
//   - A deduplicated merged slice.
func MergeStringSlices(slices ...[]string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, s := range slices {
		for _, v := range s {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// MergeInt64Slices merges multiple int64 slices, removing duplicates.
//
// Parameters:
//   - slices: The int64 slices to merge.
//
// Returns:
//   - A deduplicated merged slice.
func MergeInt64Slices(slices ...[]int64) []int64 {
	seen := make(map[int64]struct{})
	var result []int64
	for _, s := range slices {
		for _, v := range s {
			if _, ok := seen[v]; !ok {
				seen[v] = struct{}{}
				result = append(result, v)
			}
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Map merge
// ---------------------------------------------------------------------------

// MergeStringMaps merges multiple string maps. Later maps override earlier
// ones for duplicate keys.
//
// Parameters:
//   - maps: The string maps to merge.
//
// Returns:
//   - A single merged map.
func MergeStringMaps(maps ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// ---------------------------------------------------------------------------
// Deduplicate by key
// ---------------------------------------------------------------------------

// DeduplicateByKey removes duplicates from a slice using a key-extraction
// function. First occurrence wins.
//
// Parameters:
//   - items: The slice of items to deduplicate.
//   - keyFn: Function that returns the dedup key for each item.
//
// Returns:
//   - A new slice with duplicates removed.
func DeduplicateByKey[T any, K comparable](items []T, keyFn func(T) K) []T {
	seen := make(map[K]struct{})
	var result []T
	for _, item := range items {
		k := keyFn(item)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
