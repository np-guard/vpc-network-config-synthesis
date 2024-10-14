/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

const (
	aclExtenalsConfig = "%s/acl_externals/config_object.json"
	aclExternalsSpec  = "%s/acl_externals/conn_spec.json"

	aclNifConfig = "%s/acl_nif/config_object.json"
	aclNifSpec   = "%s/acl_nif/conn_spec.json"

	aclProtocolsConfig = "%s/acl_protocols/config_object.json"
	aclProtocolsSpec   = "%s/acl_protocols/conn_spec.json"

	aclSubnetCidrSegmentsConfig = "%s/acl_subnet_cidr_segments/config_object.json"
	aclSubnetCidrSegmentsSpec   = "%s/acl_subnet_cidr_segments/conn_spec.json"

	aclTesting5Config = "%s/acl_testing5/config_object.json"
	aclTesting5Spec   = "%s/acl_testing5/conn_spec.json"

	aclTgMultipleConfig = "%s/acl_tg_multiple/config_object.json"
	aclTgMultipleSpec   = "%s/acl_tg_multiple/conn_spec.json"

	sgProtocolsConfig = "%s/sg_protocols/config_object.json"
	sgProtocolsSpec   = "%s/sg_protocols/conn_spec.json"

	sgTesting3Config = "%s/sg_testing3/config_object.json"
	sgTesting3Spec   = "%s/sg_testing3/conn_spec.json"

	sgTgMultipleConfig = "%s/sg_tg_multiple/config_object.json"
	sgTgMultipleSpec   = "%s/sg_tg_multiple/conn_spec.json"

	tfOutputFmt = "tf"
)

func allMainTests() []testCase {
	return append(synthACLTestsList(), synthSGTestsList()...)
}

//nolint:funlen //all acl synthesis tests
func synthACLTestsList() []testCase {
	return []testCase{
		// acl externals    ## acl_testing4 config
		{
			testName: "acl_externals_json",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclExtenalsConfig,
				spec:       aclExternalsSpec,
				outputFile: "%s/acl_externals_json/nacl_expected.json",
			},
		},
		{
			testName: "acl_externals_tf",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclExtenalsConfig,
				spec:       aclExternalsSpec,
				outputFile: "%s/acl_externals_tf/nacl_expected.tf",
			},
		},

		// acl nif (scoping)    ## tg-multiple config
		{
			testName: "acl_nif_tf",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclNifConfig,
				spec:       aclNifSpec,
				outputFile: "%s/acl_nif_tf/nacl_expected.tf",
			},
		},

		// acl protocols (all output fmts)    ## tg-multiple config
		{
			testName: "acl_protocols_csv",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclProtocolsConfig,
				spec:       aclProtocolsSpec,
				outputFile: "%s/acl_protocols_csv/nacl_expected.csv",
			},
		},
		{
			testName: "acl_protocols_json",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclProtocolsConfig,
				spec:       aclProtocolsSpec,
				outputFile: "%s/acl_protocols_json/nacl_expected.json",
			},
		},
		{
			testName: "acl_protocols_md",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclProtocolsConfig,
				spec:       aclProtocolsSpec,
				outputFile: "%s/acl_protocols_md/nacl_expected.md",
			},
		},
		{
			testName: "acl_protocols_tf",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclProtocolsConfig,
				spec:       aclProtocolsSpec,
				outputFile: "%s/acl_protocols_tf/nacl_expected.tf",
			},
		},

		// acl subnet and cidr segments (bidi)    ## acl_testing5 config
		{
			testName: "acl_subnet_cidr_segments_tf",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclSubnetCidrSegmentsConfig,
				spec:       aclSubnetCidrSegmentsSpec,
				outputFile: "%s/acl_subnet_cidr_segments_tf/nacl_expected.tf",
			},
		},

		// acl testing 5 (json, json single, tf, tf single)
		{
			testName: "acl_testing5_json",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclTesting5Config,
				spec:       aclTesting5Spec,
				outputFile: "%s/acl_testing5_json/nacl_expected.json",
			},
		},
		{
			testName: "acl_testing5_json_single",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				singleacl:  true,
				config:     aclTesting5Config,
				spec:       aclTesting5Spec,
				outputFile: "%s/acl_testing5_json_single/nacl_single_expected.json",
			},
		},
		{
			testName: "acl_testing5_tf",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclTesting5Config,
				spec:       aclTesting5Spec,
				outputFile: "%s/acl_testing5_tf/nacl_expected.tf",
			},
		},
		{
			testName: "acl_testing5_tf_single",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				singleacl:  true,
				config:     aclTesting5Config,
				spec:       aclTesting5Spec,
				outputFile: "%s/acl_testing5_tf_single/nacl_single_expected.tf",
			},
		},

		// acl tg multiple (json, tf, tf separate)
		{
			testName: "acl_tg_multiple_json",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclTgMultipleConfig,
				spec:       aclTgMultipleSpec,
				outputFile: "%s/acl_tg_multiple_json/nacl_expected.json",
			},
		},
		{
			testName: "acl_tg_multiple_tf",
			args: &command{
				cmd:        synth,
				subcmd:     acl,
				config:     aclTgMultipleConfig,
				spec:       aclTgMultipleSpec,
				outputFile: "%s/acl_tg_multiple_tf/nacl_expected.tf",
			},
		},
		{
			testName: "acl_tg_multiple_tf_separate",
			args: &command{
				cmd:       synth,
				subcmd:    acl,
				config:    aclTgMultipleConfig,
				spec:      aclTgMultipleSpec,
				outputDir: "%s/acl_tg_multiple_tf_separate",
				format:    tfOutputFmt,
			},
		},
	}
}

func synthSGTestsList() []testCase {
	return []testCase{
		// sg protocols (all output fmts, externals, scoping, nif as a resource)    ## tg-multiple config
		{
			testName: "sg_protocols_csv",
			args: &command{
				cmd:        synth,
				subcmd:     sg,
				config:     sgProtocolsConfig,
				spec:       sgProtocolsSpec,
				outputFile: "%s/sg_protocols_csv/sg_expected.csv",
			},
		},
		{
			testName: "sg_protocols_json",
			args: &command{
				cmd:        synth,
				subcmd:     sg,
				config:     sgProtocolsConfig,
				spec:       sgProtocolsSpec,
				outputFile: "%s/sg_protocols_json/sg_expected.json",
			},
		},
		{
			testName: "sg_protocols_md",
			args: &command{
				cmd:        synth,
				subcmd:     sg,
				config:     sgProtocolsConfig,
				spec:       sgProtocolsSpec,
				outputFile: "%s/sg_protocols_md/sg_expected.md",
			},
		},
		{
			testName: "sg_protocols_tf",
			args: &command{
				cmd:        synth,
				subcmd:     sg,
				config:     sgProtocolsConfig,
				spec:       sgProtocolsSpec,
				outputFile: "%s/sg_protocols_tf/sg_expected.tf",
			},
		},

		// sg testing 3 (all fmts, VPEs are included)
		{
			testName: "sg_testing3_csv",
			args: &command{
				cmd:        synth,
				subcmd:     sg,
				config:     sgTesting3Config,
				spec:       sgTesting3Spec,
				outputFile: "%s/sg_testing3_csv/sg_expected.csv",
			},
		},
		{
			testName: "sg_testing3_json",
			args: &command{
				cmd:        synth,
				subcmd:     sg,
				config:     sgTesting3Config,
				spec:       sgTesting3Spec,
				outputFile: "%s/sg_testing3_json/sg_expected.json",
			},
		},
		{
			testName: "sg_testing3_md",
			args: &command{
				cmd:        synth,
				subcmd:     sg,
				config:     sgTesting3Config,
				spec:       sgTesting3Spec,
				outputFile: "%s/sg_testing3_md/sg_expected.md",
			},
		},
		{
			testName: "sg_testing3_tf",
			args: &command{
				cmd:        synth,
				subcmd:     sg,
				config:     sgTesting3Config,
				spec:       sgTesting3Spec,
				outputFile: "%s/sg_testing3_tf/sg_expected.tf",
			},
		},

		// sg tg multiple (tf separate)
		{
			testName: "sg_tg_multiple_tf_separate",
			args: &command{
				cmd:       synth,
				subcmd:    sg,
				config:    sgTgMultipleConfig,
				spec:      sgTgMultipleSpec,
				outputDir: "%s/sg_tg_multiple_tf_separate",
				format:    tfOutputFmt,
			},
		},
	}
}
