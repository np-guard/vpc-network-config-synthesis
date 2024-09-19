/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

func errorTestsList() []errorTestCase {
	return []errorTestCase{
		/*
			####### CLI ERRORS #######
		*/
		// -o and -d
		{
			testName: "-o with -d",
			command:  "string",
			err:      "ambiguous resource name",
		},

		// -f = json and -l

		// config was not supplied

		// give go.mod as spec

		/*
			####### INPUT ERRORS #######
		*/

		// two resources with the same name
		{
			testName: "string",
			command:  "string",
			err:      "ambiguous resource name",
		},

		// bad protocol

		// unexisting resource in spec

		// vpe resource in ACL generation
	}
}
