/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package optimize

import (
	"sort"

	"github.com/np-guard/models/pkg/ds"
	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"
	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

type Optimizer interface {
	// read the collection from the config object file
	ParseCollection(filename string) error

	// optimize number of SG/nACL rules
	Optimize() ir.OptimizeCollection

	// returns a slice of all vpc names. used to generate locals file
	VpcNames() []string
}

// each IPBlock is a single CIDR. The CIDRs are disjoint.
func sortPartitionsByIPAddrs[T any](p []ds.Pair[*netset.IPBlock, T]) []ds.Pair[*netset.IPBlock, T] {
	cmp := func(i, j int) bool { return p[i].Left.FirstIPAddress() < p[j].Left.FirstIPAddress() }
	sort.Slice(p, cmp)
	return p
}

func allPorts(ports *interval.CanonicalSet) bool {
	return ports.Equal(netp.AllPorts().ToSet())
}
