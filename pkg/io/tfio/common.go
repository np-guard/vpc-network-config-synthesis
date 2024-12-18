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

	"github.com/np-guard/models/pkg/interval"
	"github.com/np-guard/models/pkg/netp"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/io/tfio/tf"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/ir"
)

const (
	resourceConst = "resource"
	nameConst     = "name"
)

// Writer implements ir.Writer
type Writer struct {
	w *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: bufio.NewWriter(w)}
}

func portRange(r interval.Interval, prefix string) []tf.Argument {
	var arguments []tf.Argument
	if r.Start() != netp.MinPort {
		arguments = append(arguments, tf.Argument{Name: prefix + "_min", Value: strconv.FormatInt(r.Start(), 10)})
	}
	if r.End() != netp.MaxPort {
		arguments = append(arguments, tf.Argument{Name: prefix + "_max", Value: strconv.FormatInt(r.End(), 10)})
	}
	return arguments
}

func codeTypeArguments(ct *netp.ICMPTypeCode) []tf.Argument {
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

// Resource names must start with a letter or underscore, and may
// contain only letters, digits, underscores, and dashes.
// (https://developer.hashicorp.com/terraform/language/resources/syntax)
func verifyName(name string) error {
	pattern := "^[A-Za-z_][A-Za-z0-9_-]*$"
	ok, err := regexp.MatchString(pattern, name)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("%q should match regexp %q", name, pattern)
	}
	return nil
}
