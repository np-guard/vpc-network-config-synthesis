/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"fmt"
	"slices"

	"github.com/np-guard/vpc-network-config-synthesis/pkg/synth"
	"github.com/np-guard/vpc-network-config-synthesis/pkg/utils"
)

const (
	tgMultipleConfig               = "%s/tg_multiple/config_object.json"
	sgTesting3Config               = "%s/sg_testing3/config_object.json"
	aclTesting4Config              = "%s/acl_testing4/config_object.json"
	aclTesting5Config              = "%s/acl_testing5/config_object.json"
	optimizeACLConfig              = "%s/optimize_acl/config_object.json"
	optimizeACLAnyProtocolConfig   = "%s/optimize_acl_anyProtocol/config_object.json"
	optimizeSGProtocolsToAllConfig = "%s/optimize_sg_protocols_to_all/config_object.json"

	aclExternalsSpec           = "%s/acl_externals/conn_spec.json"
	aclNifSpec                 = "%s/acl_nif/conn_spec.json"
	aclNifInstanceSegmentsSpec = "%s/acl_nif_instance_segments/conn_spec.json"
	aclProtocolsSpec           = "%s/acl_protocols/conn_spec.json"
	aclSubnetCidrSegmentsSpec  = "%s/acl_subnet_cidr_segments/conn_spec.json"
	aclTesting5Spec            = "%s/acl_testing5/conn_spec.json"
	aclTgMultipleSpec          = "%s/acl_tg_multiple/conn_spec.json"
	aclVpeSpec                 = "%s/acl_vpe/conn_spec.json"
	sgProtocolsSpec            = "%s/sg_protocols/conn_spec.json"
	sgSegments1Spec            = "%s/sg_segments1/conn_spec.json"
	sgSegments2Spec            = "%s/sg_segments2/conn_spec.json"
	sgSegments3Spec            = "%s/sg_segments3/conn_spec.json"
	sgSegments4Spec            = "%s/sg_segments4/conn_spec.json"
	sgTesting3Spec             = "%s/sg_testing3/conn_spec.json"
	sgTgMultipleSpec           = "%s/sg_tg_multiple/conn_spec.json"

	tfOutputFmt = "tf"
	vsi1        = "test-vpc1--vsi1"
)

func allMainTests() []testCase {
	return slices.Concat(synthACLTestsList(), synthSGTestsList(), optimizeSGTestsLists(), optimizeACLTestsLists())
}

//nolint:funlen //all acl synthesis tests
func synthACLTestsList() []testCase {
	return []testCase{
		// acl externals
		{
			testName: "acl_externals_json",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     aclTesting4Config,
				spec:       aclExternalsSpec,
				outputFile: "%s/acl_externals_json/nacl_expected.json",
			},
		},
		{
			testName: "acl_externals_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     aclTesting4Config,
				spec:       aclExternalsSpec,
				outputFile: "%s/acl_externals_tf/nacl_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedACL, "test-vpc1/subnet3")),
		},

		// acl nif (scoping)
		{
			testName: "acl_nif_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     tgMultipleConfig,
				spec:       aclNifSpec,
				outputFile: "%s/acl_nif_tf/nacl_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedACL,
				"test-vpc0/subnet2, test-vpc0/subnet3, test-vpc0/subnet5, test-vpc1/subnet10, test-vpc1/subnet11, test-vpc2/subnet20, ",
				"test-vpc3/subnet30")),
		},

		// acl nif instance segments
		{
			testName: "acl_nif_instance_segments_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     tgMultipleConfig,
				spec:       aclNifInstanceSegmentsSpec,
				outputFile: "%s/acl_nif_instance_segments_tf/nacl_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedACL,
				"test-vpc0/subnet1, test-vpc0/subnet4, test-vpc0/subnet5, test-vpc1/subnet10, test-vpc1/subnet11, test-vpc2/subnet20")),
		},

		// acl protocols (all output fmts)
		{
			testName: "acl_protocols_csv",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     tgMultipleConfig,
				spec:       aclProtocolsSpec,
				outputFile: "%s/acl_protocols_csv/nacl_expected.csv",
			},
		},
		{
			testName: "acl_protocols_json",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     tgMultipleConfig,
				spec:       aclProtocolsSpec,
				outputFile: "%s/acl_protocols_json/nacl_expected.json",
			},
		},
		{
			testName: "acl_protocols_md",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     tgMultipleConfig,
				spec:       aclProtocolsSpec,
				outputFile: "%s/acl_protocols_md/nacl_expected.md",
			},
		},
		{
			testName: "acl_protocols_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     tgMultipleConfig,
				spec:       aclProtocolsSpec,
				outputFile: "%s/acl_protocols_tf/nacl_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedACL,
				"test-vpc2/subnet20, test-vpc3/subnet30")),
		},

		// acl subnet and cidr segments (bidi)
		{
			testName: "acl_subnet_cidr_segments_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     aclTesting5Config,
				spec:       aclSubnetCidrSegmentsSpec,
				outputFile: "%s/acl_subnet_cidr_segments_tf/nacl_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedACL,
				"testacl5-vpc/sub1-1, testacl5-vpc/sub3-1")),
		},

		// acl testing 5 (json, json single, tf, tf single)
		{
			testName: "acl_testing5_json",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     aclTesting5Config,
				spec:       aclTesting5Spec,
				outputFile: "%s/acl_testing5_json/nacl_expected.json",
			},
		},
		{
			testName: "acl_testing5_json_single",
			args: &command{
				cmd:        synthesis,
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
				cmd:        synthesis,
				subcmd:     acl,
				config:     aclTesting5Config,
				spec:       aclTesting5Spec,
				outputFile: "%s/acl_testing5_tf/nacl_expected.tf",
			},
			expectedWarning: utils.Ptr(""),
		},
		{
			testName: "acl_testing5_tf_single",
			args: &command{
				cmd:        synthesis,
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
				cmd:        synthesis,
				subcmd:     acl,
				config:     tgMultipleConfig,
				spec:       aclTgMultipleSpec,
				outputFile: "%s/acl_tg_multiple_json/nacl_expected.json",
			},
		},
		{
			testName: "acl_tg_multiple_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     tgMultipleConfig,
				spec:       aclTgMultipleSpec,
				outputFile: "%s/acl_tg_multiple_tf/nacl_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedACL,
				"test-vpc0/subnet1, test-vpc2/subnet20, test-vpc3/subnet30")),
		},
		{
			testName: "acl_tg_multiple_tf_separate",
			args: &command{
				cmd:       synthesis,
				subcmd:    acl,
				config:    tgMultipleConfig,
				spec:      aclTgMultipleSpec,
				outputDir: "%s/acl_tg_multiple_tf_separate",
				format:    tfOutputFmt,
			},
		},

		// acl vpe    ## sg_testing3 config
		{
			testName: "acl_vpe_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     acl,
				config:     sgTesting3Config,
				spec:       aclVpeSpec,
				outputFile: "%s/acl_vpe_tf/nacl_expected.tf",
			},
			expectedWarning: utils.Ptr(""),
		},
	}
}

//nolint:funlen // test cases
func synthSGTestsList() []testCase {
	return []testCase{
		// sg protocols (all output fmts, externals, scoping, nif as a resource)
		{
			testName: "sg_protocols_csv",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     tgMultipleConfig,
				spec:       sgProtocolsSpec,
				outputFile: "%s/sg_protocols_csv/sg_expected.csv",
			},
		},
		{
			testName: "sg_protocols_json",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     tgMultipleConfig,
				spec:       sgProtocolsSpec,
				outputFile: "%s/sg_protocols_json/sg_expected.json",
			},
		},
		{
			testName: "sg_protocols_md",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     tgMultipleConfig,
				spec:       sgProtocolsSpec,
				outputFile: "%s/sg_protocols_md/sg_expected.md",
			},
		},
		{
			testName: "sg_protocols_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     tgMultipleConfig,
				spec:       sgProtocolsSpec,
				outputFile: "%s/sg_protocols_tf/sg_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedSG,
				"test-vpc0/vsi0-subnet4, test-vpc0/vsi0-subnet5, test-vpc0/vsi1-subnet2, test-vpc0/vsi1-subnet3, ",
				"test-vpc0/vsi1-subnet4, test-vpc0/vsi1-subnet5, test-vpc1/vsi0-subnet11, test-vpc2/vsi0-subnet20, ",
				"test-vpc2/vsi2-subnet20, test-vpc3/vsi0-subnet30")),
		},

		// sg segments1 (cidrSegment -> cidrSegment)
		{
			testName: "sg_segments1_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     tgMultipleConfig,
				spec:       sgSegments1Spec,
				outputFile: "%s/sg_segments1_tf/sg_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedSG,
				"test-vpc0/vsi0-subnet2, test-vpc0/vsi0-subnet3, test-vpc0/vsi0-subnet4, test-vpc0/vsi0-subnet5, ",
				"test-vpc0/vsi1-subnet2, test-vpc0/vsi1-subnet3, test-vpc0/vsi1-subnet4, test-vpc0/vsi1-subnet5, ",
				"test-vpc1/vsi0-subnet10, test-vpc1/vsi0-subnet11, test-vpc2/vsi0-subnet20, ",
				"test-vpc2/vsi1-subnet20, test-vpc2/vsi2-subnet20, test-vpc3/vsi0-subnet30")),
		},

		// sg segments2 (instanceSegment -> cidrSegment)
		{
			testName: "sg_segments2_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     tgMultipleConfig,
				spec:       sgSegments2Spec,
				outputFile: "%s/sg_segments2_tf/sg_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedSG,
				"test-vpc0/vsi0-subnet2, test-vpc0/vsi0-subnet3, test-vpc0/vsi0-subnet4, test-vpc0/vsi1-subnet2, ",
				"test-vpc0/vsi1-subnet3, test-vpc0/vsi1-subnet4, test-vpc0/vsi1-subnet5, test-vpc1/vsi0-subnet10, ",
				"test-vpc2/vsi0-subnet20, test-vpc2/vsi1-subnet20, test-vpc2/vsi2-subnet20, test-vpc3/vsi0-subnet30")),
		},

		// sg segments3 (subnetSegment -> nifSegment)
		{
			testName: "sg_segments3_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     tgMultipleConfig,
				spec:       sgSegments3Spec,
				outputFile: "%s/sg_segments3_tf/sg_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedSG,
				"test-vpc0/vsi0-subnet0, test-vpc0/vsi0-subnet1, test-vpc0/vsi0-subnet2, test-vpc0/vsi0-subnet3, ",
				"test-vpc0/vsi0-subnet5, test-vpc0/vsi1-subnet0, test-vpc0/vsi1-subnet1, test-vpc0/vsi1-subnet2, ",
				"test-vpc0/vsi1-subnet3, test-vpc0/vsi1-subnet4, test-vpc1/vsi0-subnet11, test-vpc3/vsi0-subnet30")),
		},

		// sg segments4 (vpeSegment -> instanceSegment)
		{
			testName: "sg_segments4_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     sgTesting3Config,
				spec:       sgSegments4Spec,
				outputFile: "%s/sg_segments4_tf/sg_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedSG, "test-vpc/opa, test-vpc/proxy")),
		},

		// sg testing 3 (all fmts, VPEs are included)
		{
			testName: "sg_testing3_csv",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     sgTesting3Config,
				spec:       sgTesting3Spec,
				outputFile: "%s/sg_testing3_csv/sg_expected.csv",
			},
		},
		{
			testName: "sg_testing3_json",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     sgTesting3Config,
				spec:       sgTesting3Spec,
				outputFile: "%s/sg_testing3_json/sg_expected.json",
			},
		},
		{
			testName: "sg_testing3_md",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     sgTesting3Config,
				spec:       sgTesting3Spec,
				outputFile: "%s/sg_testing3_md/sg_expected.md",
			},
		},
		{
			testName: "sg_testing3_tf",
			args: &command{
				cmd:        synthesis,
				subcmd:     sg,
				config:     sgTesting3Config,
				spec:       sgTesting3Spec,
				outputFile: "%s/sg_testing3_tf/sg_expected.tf",
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedSG, "test-vpc/appdata-endpoint-gateway")),
		},

		// sg tg multiple (tf separate)
		{
			testName: "sg_tg_multiple_tf_separate",
			args: &command{
				cmd:       synthesis,
				subcmd:    sg,
				config:    tgMultipleConfig,
				spec:      sgTgMultipleSpec,
				outputDir: "%s/sg_tg_multiple_tf_separate",
				format:    tfOutputFmt,
			},
			expectedWarning: utils.Ptr(fmt.Sprint(synth.WarningUnspecifiedSG,
				"test-vpc0/vsi0-subnet1, test-vpc0/vsi0-subnet2, test-vpc0/vsi0-subnet3, test-vpc0/vsi0-subnet4, ",
				"test-vpc0/vsi0-subnet5, test-vpc0/vsi1-subnet0, test-vpc0/vsi1-subnet1, test-vpc0/vsi1-subnet2,",
				" test-vpc0/vsi1-subnet3, test-vpc0/vsi1-subnet5, test-vpc2/vsi1-subnet20, test-vpc3/vsi0-subnet30")),
		},
	}
}

// Note1: spec files in data folder are used to create the config object files (acl_testing4 config)
// Note2: each data folder has a details.txt file with the test explanation
func optimizeSGTestsLists() []testCase {
	return []testCase{
		// optimize_sg_protocols_to_all tests also test SG rules with local values different from 0.0.0.0/0
		{
			testName: "optimize_sg_protocols_to_all_tf",
			args: &command{
				cmd:        optimize,
				subcmd:     sg,
				config:     optimizeSGProtocolsToAllConfig,
				outputFile: "%s/optimize_sg_protocols_to_all_tf/sg_expected.tf",
			},
		},
		{
			testName: "optimize_sg_protocols_to_all_csv",
			args: &command{
				cmd:        optimize,
				subcmd:     sg,
				config:     optimizeSGProtocolsToAllConfig,
				outputFile: "%s/optimize_sg_protocols_to_all_csv/sg_expected.csv",
			},
		},
		{
			testName: "optimize_sg_protocols_to_all_json",
			args: &command{
				cmd:        optimize,
				subcmd:     sg,
				config:     optimizeSGProtocolsToAllConfig,
				outputFile: "%s/optimize_sg_protocols_to_all_json/sg_expected.json",
			},
		},
		{
			testName: "optimize_sg_protocols_to_all_md",
			args: &command{
				cmd:        optimize,
				subcmd:     sg,
				config:     optimizeSGProtocolsToAllConfig,
				outputFile: "%s/optimize_sg_protocols_to_all_md/sg_expected.md",
			},
		},
		{
			testName: "optimize_sg_redundant",
			args: &command{
				cmd:        optimize,
				subcmd:     sg,
				config:     "%s/optimize_sg_redundant/config_object.json",
				outputFile: "%s/optimize_sg_redundant/sg_expected.tf",
			},
		},
		{
			testName: "optimize_sg_t",
			args: &command{
				cmd:          optimize,
				subcmd:       sg,
				config:       "%s/optimize_sg_t/config_object.json",
				outputFile:   "%s/optimize_sg_t/sg_expected.tf",
				firewallName: vsi1,
			},
		},
		{
			testName: "optimize_sg_t_all",
			args: &command{
				cmd:          optimize,
				subcmd:       sg,
				config:       "%s/optimize_sg_t_all/config_object.json",
				outputFile:   "%s/optimize_sg_t_all/sg_expected.tf",
				firewallName: vsi1,
			},
		},
	}
}

func optimizeACLTestsLists() []testCase {
	return []testCase{
		{
			testName: "optimize_acl_csv",
			args: &command{
				cmd:          optimize,
				subcmd:       acl,
				config:       optimizeACLConfig,
				outputFile:   "%s/optimize_acl_csv/nacl_expected.csv",
				firewallName: "testacl5-vpc--sub1-2",
			},
		},
		{
			testName: "optimize_acl_md",
			args: &command{
				cmd:        optimize,
				subcmd:     acl,
				config:     optimizeACLConfig,
				outputFile: "%s/optimize_acl_md/nacl_expected.md",
			},
		},
		{
			testName: "optimize_acl_tf",
			args: &command{
				cmd:        optimize,
				subcmd:     acl,
				config:     optimizeACLConfig,
				outputFile: "%s/optimize_acl_tf/nacl_expected.tf",
			},
		},
		{
			testName: "optimize_acl_json",
			args: &command{
				cmd:        optimize,
				subcmd:     acl,
				config:     optimizeACLConfig,
				outputFile: "%s/optimize_acl_json/nacl_expected.json",
			},
		},
		{
			testName: "optimize_acl_anyProtocol_tf",
			args: &command{
				cmd:        optimize,
				subcmd:     acl,
				config:     optimizeACLAnyProtocolConfig,
				outputFile: "%s/optimize_acl_anyProtocol_tf/nacl_expected.tf",
			},
		},
	}
}
