/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import "github.com/np-guard/models/pkg/netset"

// temporary file, should be implement in models repo

// given a<b disjoint, returns true if a and b are touching
func touching(a *netset.IPBlock, b *netset.IPBlock) bool {
	return len(a.Union(b).Split()) == 1
}
