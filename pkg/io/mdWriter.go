/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package io

import (
	"bufio"
	"io"
	"slices"
	"strings"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const (
	leftAlign = " :--- "
	separator = " | "
)

// MDWriter implements ir.Writer
type MDWriter struct {
	w *bufio.Writer
}

func NewMDWriter(w io.Writer) *MDWriter {
	return &MDWriter{w: bufio.NewWriter(w)}
}

func (w *MDWriter) WriteSG(collection *ir.SGCollection, vpc string, _ bool) error {
	sgTable, err := WriteSG(collection, vpc)
	if err != nil {
		return err
	}
	sgHeader := makeSGHeader()
	return w.writeAll(slices.Concat(sgHeader, addAligns(len(sgHeader[0])), sgTable))
}

func (w *MDWriter) WriteACL(collection *ir.ACLCollection, vpc string, _ bool) error {
	aclTable, err := WriteACL(collection, vpc)
	if err != nil {
		return err
	}
	aclHeader := makeACLHeader()
	return w.writeAll(slices.Concat(aclHeader, addAligns(len(aclHeader[0])), aclTable))
}

func (w *MDWriter) writeAll(rows [][]string) error {
	for _, row := range rows {
		finalString := separator + strings.Join(row, separator) + separator + "\n"
		if _, err := w.w.WriteString(finalString); err != nil {
			return err
		}
	}
	return w.w.Flush()
}

func addAligns(n int) [][]string {
	res := make([]string, n)
	for i := range n {
		res[i] = leftAlign
	}
	return [][]string{res}
}
