/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

//nolint:lll // commands can be long
func errorTestsList() []errorTestCase {
	return []errorTestCase{
		/*  ############################  */
		/*	####### CLI ERRORS #######  */
		/*  ############################  */
		// -o and -d
		{
			testName: "-o with -d",
			command:  "../bin/vpcgen synth acl -c data_errors/cli/config_object.json -s data_errors/cli/conn_spec.json -o data_errors/cli/nacl_expected.tf -d data_errors/cli",
			err:      "ambiguous resource name", // check
		},

		// -f = json and -l
		{
			testName: "locals jon fmt",
			command:  "../bin/vpcgen synth acl -c -o data_errors/cli/config_object.json -s data_errors/cli/conn_spec.json -o data_errors/cli/nacl_expected.json -l",
			err:      "--locals flag requires setting the output format to tf",
		},

		// json fmt with -d
		{
			testName: "json separate",
			command:  "../bin/vpcgen synth acl -c -o data_errors/cli/config_object.json -s data_errors/cli/conn_spec.json -d data_errors/cli -f json ",
			err:      "--locals flag requires setting the output format to tf", // check
		},

		// config was not supplied
		{
			testName: "no config file",
			command:  "../bin/vpcgen synth acl -c data_errors/cli/config_object.json -s data_errors/cli/conn_spec.json -o data_errors/cli/nacl_expected.tf",
			err:      "--locals flag requires setting the output format to tf", // check
		},

		// give go.mod as spec
		{
			testName: "bad spec file",
			command:  "../bin/vpcgen synth acl -c data_errors/cli/config_object.json -s ../go.mod -o data_errors/cli/nacl_expected.tf",
			err:      "--locals flag requires setting the output format to tf", // check
		},

		// unknown subcmd
		{
			testName: "unknown subcmd",
			command:  "../bin/vpcgen synth acl -c data_errors/cli/config_object.json -s conn_spec.json -o data_errors/cli/nacl_expected.tf",
			err:      "--locals flag requires setting the output format to tf", // check
		},

		/*  ############################  */
		/*	####### INPUT ERRORS #######  */
		/*  ############################  */

		// two resources with the same name
		{
			testName: "ambiguous resource name",
			command:  "../bin/vpcgen synth acl -c data_errors/ambiguous/config_object.json -s data_errors/ambiguous/conn_spec.json -o data_errors/ambiguous/nacl_expected.tf",
			err:      "ambiguous resource name",
		},

		// bad protocol
		{
			testName: "bad protocol",
			command:  "../bin/vpcgen synth acl -c %s/bad_protocol/config_object.json -s %s/bad_protocol/conn_spec.json -o %s/bad_protocol/nacl_expected.tf",
			err:      "ambiguous resource name", // check
		},

		// unknown resource in spec
		{
			testName: "unknown resource",
			command:  "../bin/vpcgen synth acl -c %s/unknown_resource/config_object.json -s %s/unknown_resource/conn_spec.json -o %s/unknown_resource/nacl_expected.tf",
			err:      "ambiguous resource name", // check
		},

		// vpe resource in ACL generation
		{
			testName: "vpe_acl",
			command:  "../bin/vpcgen synth acl -c %s/vpe_acl/config_object.json -s %s/vpe_acl/conn_spec.json -o %s/vpe_acl/nacl_expected.tf",
			err:      "ambiguous resource name", // check
		},
	}
}
