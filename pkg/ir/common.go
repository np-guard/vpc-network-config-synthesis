/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

type (
	Direction string

	SynthCollection interface {
		WriteSynth(writer SynthWriter, vpc string) error
	}

	OptimizeCollection interface {
		WriteOptimize(writer OptimizeWriter) error
	}

	SynthWriter interface {
		SynthACLWriter
		SynthSGWriter
	}

	OptimizeWriter interface {
		OptimizeACLWriter
		OptimizeSGWriter
	}
)

const (
	Outbound Direction = "outbound"
	Inbound  Direction = "inbound"
)
