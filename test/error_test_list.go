/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

//nolint:lll // commands can be long
func errorTestsList() []errorTestCase {
	return []errorTestCase{
		/*  ############################  */
		/*	####### CLI ERRORS #########  */
		/*  ############################  */
		// -o and -d
		{
			testName: "-o with -d",
			command:  "../bin/vpcgen synth acl -c data_errors/cli/config_object.json -s data_errors/cli/conn_spec.json -o data_errors/cli/nacl_expected.tf -d data_errors/cli",
			err:      "if any flags in the group [output-file output-dir] are set none of the others can be; [output-dir output-file] were all set",
		},

		// -f = json and -l
		{
			testName: "locals json fmt",
			command:  "../bin/vpcgen synth acl -c data_errors/cli/config_object.json -s data_errors/cli/conn_spec.json -o data_errors/cli/nacl_expected.json -l",
			err:      "--locals flag requires setting the output format to tf",
		},

		// json fmt with -d
		{
			testName: "json separate",
			command:  "../bin/vpcgen synth acl -c data_errors/cli/config_object.json -s data_errors/cli/conn_spec.json -d data_errors/cli -f json",
			err:      "-d cannot be used with format json",
		},

		// config was not supplied
		{
			testName: "no config file",
			command:  "../bin/vpcgen synth acl -s data_errors/cli/conn_spec.json -o data_errors/cli/nacl_expected.tf",
			err:      "required flag(s) \"config\" not set",
		},

		// give go.mod as spec
		{
			testName: "bad spec file",
			command:  "../bin/vpcgen synth acl -c data_errors/cli/config_object.json -s ../go.mod -o data_errors/cli/nacl_expected.tf",
			err:      "could not parse connectivity file",
		},

		// unknown subcmd
		{
			testName: "unknown subcmd",
			command:  "../bin/vpcgen pop acl -c data_errors/cli/config_object.json -s conn_spec.json -o data_errors/cli/nacl_expected.tf",
			err:      "unknown command \"pop\" for \"vpcgen\"",
		},

		/*  ############################  */
		/*	####### INPUT ERRORS #######  */
		/*  ############################  */

		// two resources with the same name
		{
			testName: "ambiguous resource name",
			command:  "../bin/vpcgen synth acl -c data_errors/ambiguous/config_object.json -s data_errors/ambiguous/conn_spec.json -o data_errors/ambiguous/nacl_expected.tf",
			err:      "ambiguous resource name: subnet0",
		},

		// bad protocol
		{
			testName: "bad protocol",
			command:  "../bin/vpcgen synth acl -c data_errors/bad_protocol/config_object.json -s data_errors/bad_protocol/conn_spec.json -o data_errors/bad_protocol/nacl_expected.tf",
			err:      "could not parse connectivity file data_errors/bad_protocol/conn_spec.json: invalid protocol type \"ALOHA\"",
		},

		// unknown resource in spec
		{
			testName: "unknown resource",
			command:  "../bin/vpcgen synth acl -c data_errors/unknown_resource/config_object.json -s data_errors/unknown_resource/conn_spec.json -o data_errors/unknown_resource/nacl_expected.tf",
			err:      "unknown resource name subnet35 (resource type: \"subnet\")",
		},

		// vpe resource in ACL generation
		{
			testName: "vpe acl",
			command:  "../bin/vpcgen synth acl -c data_errors/vpe_acl/config_object.json -s data_errors/vpe_acl/conn_spec.json -o data_errors/vpe_acl/nacl_expected.tf",
			err:      "both source and destination are external for connection",
		},

		// impossible resource type
		{
			testName: "impossible resource type",
			command:  "../bin/vpcgen synth acl -c data_errors/impossible_resource_type/config_object.json -s data_errors/impossible_resource_type/conn_spec.json -o data_errors/impossible_resource_type/nacl_expected.tf",
			err:      "could not parse connectivity file data_errors/impossible_resource_type/conn_spec.json: invalid value (expected one of []interface {}{\"external\", \"segment\", \"subnet\", \"instance\", \"nif\", \"cidr\", \"vpe\"}): \"policydb-endpoint-gateway\"",
		},
	}
}
