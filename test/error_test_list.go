/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

const (
	cliConfig  = "%s/cli/config_object.json"
	cliSpec    = "%s/cli/conn_spec.json"
	outputPath = "%s/cli/nacl_expected.json"
)

//nolint:funlen //commands can be long
func errorTestsList() []testCase {
	return []testCase{
		/*  ############################  */
		/*	####### CLI ERRORS #########  */
		/*  ############################  */
		// -f = json and -l
		{
			testName:    "locals json fmt",
			expectedErr: "--locals flag requires setting the output format to tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     cliConfig,
				spec:       cliSpec,
				outputFile: outputPath,
				locals:     true,
			},
		},

		// json fmt with -d
		{
			testName:    "json separate",
			expectedErr: "-d cannot be used with format json",
			args: &command{
				cmd:       synthesis,
				subcmd:    acl,
				config:    cliConfig,
				spec:      cliSpec,
				outputDir: outputPath,
				format:    "json",
			},
		},

		// config was not supplied
		{
			testName:    "no config file",
			expectedErr: "required flag(s) \"config\" not set",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				spec:       cliSpec,
				outputFile: outputPath,
			},
		},

		// give config as spec
		{
			testName:    "bad spec file",
			expectedErr: "could not parse connectivity file",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     cliConfig,
				spec:       cliConfig,
				outputFile: outputPath,
			},
		},

		// unknown subcmd
		{
			testName:    "unknown subcmd",
			expectedErr: "unknown command \"pop\" for \"vpcgen\"",
			args: &command{
				cmd:        "pop",
				subcmd:     acl,
				config:     cliConfig,
				spec:       cliSpec,
				outputFile: outputPath,
			},
		},

		/*  ############################  */
		/*	####### INPUT ERRORS #######  */
		/*  ############################  */

		// two resources with the same name
		{
			testName:    "ambiguous resource name",
			expectedErr: "ambiguous resource name: subnet0",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     "%s/ambiguous/config_object.json",
				spec:       "%s/ambiguous/conn_spec.json",
				outputFile: outputPath,
			},
		},

		// bad protocol
		{
			testName:    "bad protocol",
			expectedErr: "could not parse connectivity file data_for_testing_errors/bad_protocol/conn_spec.json: invalid protocol type \"ALOHA\"",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     "%s/bad_protocol/config_object.json",
				spec:       "%s/bad_protocol/conn_spec.json",
				outputFile: outputPath,
			},
		},

		// external src and dst
		{
			testName:    "externals src and dst",
			expectedErr: "both source (dns) and destination (public internet) are external in required connection",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     "%s/externals/config_object.json",
				spec:       "%s/externals/conn_spec.json",
				outputFile: outputPath,
			},
		},

		// unknown resource in spec
		{
			testName:    "unknown resource",
			expectedErr: "unknown resource name subnet35 (resource type: \"subnet\")",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     "%s/unknown_resource/config_object.json",
				spec:       "%s/unknown_resource/conn_spec.json",
				outputFile: outputPath,
			},
		},

		// impossible resource type
		{
			testName: "impossible resource type",
			expectedErr: "could not parse connectivity file data_for_testing_errors/impossible_resource_type/conn_spec.json: " +
				"invalid value (expected one of []interface {}{\"external\", \"segment\", \"subnet\"," +
				" \"instance\", \"nif\", \"cidr\", \"vpe\"}): \"policydb-endpoint-gateway\"",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     "%s/impossible_resource_type/config_object.json",
				spec:       "%s/impossible_resource_type/conn_spec.json",
				outputFile: outputPath,
			},
		},
	}
}
