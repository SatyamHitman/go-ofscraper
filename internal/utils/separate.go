// =============================================================================
// FILE: internal/utils/separate.go
// PURPOSE: Separation/partitioning utilities. Provides functions to split
//          slices and data into groups based on predicates or keys. Used for
//          separating media by type, posts by source, etc. Ports Python
//          utils/separate.py.
// =============================================================================

package utils

// ---------------------------------------------------------------------------
// Partition
// ---------------------------------------------------------------------------

// Partition splits a slice into two groups based on a predicate. Items where
// the predicate returns true go into the first slice, others into the second.
//
// Parameters:
//   - items: The slice to partition.
//   - pred: Predicate function.
//
// Returns:
//   - (matching, notMatching) slices.
func Partition[T any](items []T, pred func(T) bool) ([]T, []T) {
	var yes, no []T
	for _, item := range items {
		if pred(item) {
			yes = append(yes, item)
		} else {
			no = append(no, item)
		}
	}
	return yes, no
}

// GroupBy groups slice items by a key-extraction function.
//
// Parameters:
//   - items: The slice to group.
//   - keyFn: Returns the grouping key for each item.
//
// Returns:
//   - A map from key to slice of items with that key.
func GroupBy[T any, K comparable](items []T, keyFn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range items {
		k := keyFn(item)
		result[k] = append(result[k], item)
	}
	return result
}

// ChunkSlice splits a slice into chunks of at most the given size.
//
// Parameters:
//   - items: The slice to chunk.
//   - size: Maximum chunk size. Must be > 0.
//
// Returns:
//   - A slice of chunks.
func ChunkSlice[T any](items []T, size int) [][]T {
	if size <= 0 {
		return [][]T{items}
	}
	var chunks [][]T
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[i:end])
	}
	return chunks
}

// FlatMap applies a function to each item and flattens the result slices.
//
// Parameters:
//   - items: Input slice.
//   - fn: Function that returns a slice for each item.
//
// Returns:
//   - Flattened result slice.
func FlatMap[T any, U any](items []T, fn func(T) []U) []U {
	var result []U
	for _, item := range items {
		result = append(result, fn(item)...)
	}
	return result
}
