/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package io

import (
	"encoding/csv"
	"io"
	"slices"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// CSVWriter implements ir.Writer
type CSVWriter struct {
	w *csv.Writer
}

func NewCSVWriter(w io.Writer) *CSVWriter {
	return &CSVWriter{w: csv.NewWriter(w)}
}

func (w *CSVWriter) WriteSG(collection *ir.SGCollection, vpc string) error {
	sgTable, err := WriteSG(collection, vpc)
	if err != nil {
		return err
	}
	return w.w.WriteAll(slices.Concat(makeSGHeader(), sgTable))
}

func (w *CSVWriter) WriteACL(collection *ir.ACLCollection, vpc string) error {
	aclTable, err := WriteACL(collection, vpc)
	if err != nil {
		return err
	}
	return w.w.WriteAll(slices.Concat(makeACLHeader(), aclTable))
}
