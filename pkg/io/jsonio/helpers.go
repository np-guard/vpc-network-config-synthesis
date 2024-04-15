/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package jsonio

import "fmt"

func enterField(field, ctx string, j map[string]interface{}) (resCtx string, res map[string]interface{}) {
	resCtx = ctx + "." + field
	res, ok := j[field].(map[string]interface{})
	if !ok {
		res = nil
	}
	return
}

func enterArray[T any](i int, ctx string, j []interface{}) (resCtx string, res T) {
	resCtx = ctx + "." + fmt.Sprintf("[%v]", i)
	res, ok := j[i].(T)
	if !ok {
		res = *new(T)
	}
	return
}

type jsonType interface {
	string | bool | float64 | *int | []string | map[string]interface{} | []interface{}
}

func enter[T jsonType](field, ctx string, j map[string]interface{}) (resCtx string, res T) {
	resCtx = ctx + "." + field
	res, ok := j[field].(T)
	if !ok {
		// TODO: use default argument
		res = *new(T)
	}
	return
}

func enterInt(field, ctx string, j map[string]interface{}) (resCtx string, result int) {
	resCtx, tmp := enter[float64](field, ctx, j)
	return resCtx, int(tmp)
}
