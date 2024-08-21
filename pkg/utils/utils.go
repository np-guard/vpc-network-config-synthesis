/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"maps"
	"slices"
	"sort"
)

func Ptr[T any](t T) *T {
	return &t
}

func MapKeys[T comparable, V any](m map[T]V) []T {
	keys := maps.Keys(m)
	result := make([]T, 0)
	for key := range keys {
		result = append(result, key)
	}
	return result
}

func SortedMapKeys[T ~string, V any](m map[T]V) []T {
	keys := MapKeys(m)
	slices.Sort(keys)
	return keys
}

func SortedAllInnerMapsKeys[T ~string, V any](m map[string]map[T]V) []T {
	keys := make([]T, 0)
	for _, vpc := range m {
		keys = append(keys, MapKeys(vpc)...)
	}
	cmp := func(i, j int) bool { return keys[i] < keys[j] }
	sort.Slice(keys, cmp)

	return keys
}

func SortedInnerMapKeys[T ~string, V any](m map[string]map[T]V, key string) []T {
	keys := MapKeys(m[key])
	cmp := func(i, j int) bool { return keys[i] < keys[j] }
	sort.Slice(keys, cmp)

	return keys
}
