/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"cmp"
	"maps"
	"slices"
)

func Ptr[T any](t T) *T {
	return &t
}

func MapKeys[T comparable, V any](m map[T]V) []T {
	return slices.Collect(maps.Keys(m))
}

func SortedMapKeys[T cmp.Ordered, V any](m map[T]V) []T {
	return slices.Sorted(maps.Keys(m))
}

func SortedAllInnerMapsKeys[T, K cmp.Ordered, V any](m map[K]map[T]V) []T {
	keys := make([]T, 0)
	for _, vpc := range m {
		keys = append(keys, MapKeys(vpc)...)
	}
	slices.Sort(keys)
	return keys
}

func TrueKeyValues[T cmp.Ordered](m map[T]bool) []T {
	keys := make([]T, 0)
	for _, k := range SortedMapKeys(m) {
		if m[k] {
			keys = append(keys, k)
		}
	}
	return keys
}
