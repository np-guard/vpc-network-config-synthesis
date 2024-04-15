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

func SortedKeys[T ~string, V any](m map[T]V) []T {
	keys := maps.Keys(m)
	cmp := func(i, j int) bool { return keys[i] < keys[j] }
	sort.Slice(keys, cmp)
	return keys
}
