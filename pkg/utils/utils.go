/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"cmp"
	"maps"
	"slices"

	"github.com/np-guard/models/pkg/ds"
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

// m1 does not have any key of m2
func MergeSetMaps[T comparable, K ds.Set[K]](m1, m2 map[T]K) map[T]K {
	for key, val := range m2 {
		m1[key] = val.Copy()
	}
	return m1
}
