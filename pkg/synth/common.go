/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package synth

import (
	"fmt"

	"github.com/np-guard/models/pkg/netset"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

type (
	Synthesizer interface {
		Synth() (ir.Collection, error)
	}

	namedAddrs struct {
		Name  ir.ID
		Addrs *netset.IPBlock
	}

	explanation struct {
		isResponse       bool
		internal         bool
		connectionOrigin fmt.Stringer
		protocolOrigin   fmt.Stringer
	}
)

func (e explanation) response() explanation {
	e.isResponse = true
	return e
}

func (e explanation) String() string {
	locality := "External"
	if e.internal {
		locality = "Internal"
	}
	result := fmt.Sprintf("%v; %v", e.connectionOrigin, e.protocolOrigin)
	if e.isResponse {
		result = fmt.Sprintf("response to %v", result)
	}
	result = fmt.Sprintf("%v. %v", locality, result)
	return result
}

func internalConn(conn *ir.Connection) (internalSrc, internalDst, internal bool) {
	internalSrc = conn.Src.Type != ir.ResourceTypeExternal
	internalDst = conn.Dst.Type != ir.ResourceTypeExternal
	internal = internalSrc && internalDst
	return
}
