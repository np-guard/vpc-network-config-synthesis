/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package ir

type Direction string

const (
	Outbound Direction = "outbound"
	Inbound  Direction = "inbound"
)

type Writer interface {
	ACLWriter
	SGWriter
}

type Collection interface {
	Write(writer Writer, vpc string) error
}
