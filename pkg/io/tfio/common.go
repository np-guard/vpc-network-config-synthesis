/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package tfio implements output of ACLs and security groups in terraform format
package tfio

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

// Writer implements ir.Writer
type Writer struct {
	w *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

func portRange(r ir.PortRange, prefix string) []tf.Argument {
	var arguments []tf.Argument
	if r.Min != ir.DefaultMinPort {
		arguments = append(arguments, tf.Argument{Name: prefix + "_min", Value: strconv.Itoa(r.Min)})
	}
	if r.Max != ir.DefaultMaxPort {
		arguments = append(arguments, tf.Argument{Name: prefix + "_max", Value: strconv.Itoa(r.Max)})
	}
	return arguments
}

func codeTypeArguments(ct *ir.ICMPCodeType) []tf.Argument {
	var arguments []tf.Argument
	if ct != nil {
		arguments = append(arguments, tf.Argument{Name: "type", Value: strconv.Itoa(ct.Type)})
		if ct.Code != nil {
			arguments = append(arguments, tf.Argument{Name: "code", Value: strconv.Itoa(*ct.Code)})
		}
	}
	return arguments
}

func quote(s string) string {
	return fmt.Sprintf("%q", s)
}

func action(a ir.Action) string {
	return string(a)
}

func direction(d ir.Direction) string {
	return string(d)
}

func verifyName(name string) error {
	pattern := "^[A-Za-z_][A-Za-z0-9_-]*$"
	ok, err := regexp.MatchString(pattern, name)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("\"name\" should match regexp %q", pattern)
	}
	return nil
}
