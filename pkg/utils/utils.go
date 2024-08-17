/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"slices"
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

func SortedMapKeys[T ~string, V any](m map[T]V) []T {
	keys := maps.Keys(m)
	slices.Sort(keys)

	return keys
}

// GetProperty returns pointer p if it is valid, else it returns the provided default value
// used to get min/max port or icmp type
func GetProperty(p *int64, defaultP int64) int64 {
	if p == nil {
		return defaultP
	}
	return *p
}

func Int64PointerToIntPointer(v *int64) *int {
	if v == nil {
		return nil
	}
	return Ptr(int(*v))
}
