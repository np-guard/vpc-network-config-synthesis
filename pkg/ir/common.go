/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

type (
	Direction string

	Collection interface {
		Write(writer Writer, vpc string, isSynth bool) error
		VpcNames() []string
	}

	Writer interface {
		ACLWriter
		SGWriter
	}
)

const (
	Outbound Direction = "outbound"
	Inbound  Direction = "inbound"
)
