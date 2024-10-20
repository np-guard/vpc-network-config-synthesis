/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package io

import (
	"bufio"
	"io"
	"strings"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const (
	sgColsNum  = 7
	aclColsNum = 10

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
	return w.writeAll(append(append(SGHeader(), addAlighns(sgColsNum)...), sgTable...))
}

func (w *MDWriter) WriteACL(collection *ir.ACLCollection, vpc string, _ bool) error {
	aclTable, err := WriteACL(collection, vpc)
	if err != nil {
		return err
	}
	return w.writeAll(append(append(ACLHeader(), addAlighns(aclColsNum)...), aclTable...))
}

func (w *MDWriter) writeAll(rows [][]string) error {
	for _, row := range rows {
		if _, err := w.w.WriteString(separator); err != nil {
			return err
		}
		if _, err := w.w.WriteString(strings.Join(row, separator)); err != nil {
			return err
		}
		if _, err := w.w.WriteString(separator + "\n"); err != nil {
			return err
		}
	}
	w.w.Flush()
	return nil
}

func addAlighns(n int) [][]string {
	res := make([]string, n)
	for i := range n {
		res[i] = leftAlign
	}
	return [][]string{res}
}
