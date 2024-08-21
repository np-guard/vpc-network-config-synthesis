/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"maps"
	"slices"
)

func Ptr[T any](t T) *T {
	return &t
}

func MapKeys[T comparable, V any](m map[T]V) []T {
	return slices.Collect(maps.Keys(m))
}

func SortedMapKeys[T ~string, V any](m map[T]V) []T {
	return slices.Sorted(maps.Keys(m))
}

func SortedAllInnerMapsKeys[T ~string, V any](m map[string]map[T]V) []T {
	keys := make([]T, 0)
	for _, vpc := range m {
		keys = append(keys, MapKeys(vpc)...)
	}
	slices.Sort(keys)
	return keys
}
