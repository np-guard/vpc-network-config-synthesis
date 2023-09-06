package utils

import (
	"sort"

	"golang.org/x/exp/maps"
)

func SortedKeys[T ~string, V any](m map[T]V) []T {
	keys := maps.Keys(m)
	cmp := func(i, j int) bool { return keys[i] < keys[j] }
	sort.Slice(keys, cmp)
	return keys
}
