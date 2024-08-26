/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

type (
	Direction string

	SynthWriter interface {
		SynthACLWriter
		SynthSGWriter
	}
)

const (
	Outbound Direction = "outbound"
	Inbound  Direction = "inbound"
)

type Collection interface {
	Write(writer SynthWriter, vpc string) error
}
