/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"sort"

	"golang.org/x/exp/maps"
)

func Ptr[T any](t T) *T {
	return &t
}

func SortedKeys[T ~string, V any](m map[string]map[T]V) []T {
	keys := make([]T, 0)
	for _, vpc := range m {
		keys = append(keys, maps.Keys(vpc)...)
	}
	cmp := func(i, j int) bool { return keys[i] < keys[j] }
	sort.Slice(keys, cmp)

	return keys
}

func SortedValuesInKey[T ~string, V any](m map[string]map[T]V, key string) []T {
	keys := maps.Keys(m[key])
	cmp := func(i, j int) bool { return keys[i] < keys[j] }
	sort.Slice(keys, cmp)

	return keys
}
