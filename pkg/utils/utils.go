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
		keys = slices.Concat(keys, MapKeys(vpc))
	}
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

func TrueKeyValues[T cmp.Ordered](m map[T]bool) []T {
	keys := make([]T, 0)
	for _, k := range SortedMapKeys(m) {
		if m[k] {
			keys = append(keys, k)
		}
	}
	return keys
}
