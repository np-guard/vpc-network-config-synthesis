/*
Copyright 2023- IBM Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package test

func allMainTests() []mainTestCase {
	return append(synthACLTestsList(), synthSGTestsList()...)
}

//nolint:lll // commands can be long
func synthACLTestsList() []mainTestCase {
	return []mainTestCase{
		// acl externals    ## acl_testing4 config
		{
			testName: "acl_externals_json",
			command:  "../bin/vpcgen synth acl -c %s/acl_externals/config_object.json -s %s/acl_externals/conn_spec.json -o %s/acl_externals_json/nacl_expected.json",
		},
		{
			testName: "acl_externals_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_externals/config_object.json -s %s/acl_externals/conn_spec.json -o %s/acl_externals_tf/nacl_expected.tf",
		},

		// acl nif (scoping)    ## tg-multiple config
		{
			testName: "acl_nif_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_nif/config_object.json -s %s/acl_nif/conn_spec.json -o %s/acl_nif_tf/nacl_expected.tf",
		},

		// acl protocols (all output fmts)    ## tg-multiple config
		{
			testName: "acl_protocols_csv",
			command:  "../bin/vpcgen synth acl -c %s/acl_protocols/config_object.json -s %s/acl_protocols/conn_spec.json -o %s/acl_protocols_csv/nacl_expected.csv",
		},
		{
			testName: "acl_protocols_json",
			command:  "../bin/vpcgen synth acl -c %s/acl_protocols/config_object.json -s %s/acl_protocols/conn_spec.json -o %s/acl_protocols_json/nacl_expected.json",
		},
		{
			testName: "acl_protocols_md",
			command:  "../bin/vpcgen synth acl -c %s/acl_protocols/config_object.json -s %s/acl_protocols/conn_spec.json -o %s/acl_protocols_md/nacl_expected.md",
		},
		{
			testName: "acl_protocols_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_protocols/config_object.json -s %s/acl_protocols/conn_spec.json -o %s/acl_protocols_tf/nacl_expected.tf",
		},

		// acl segments (bidi)    ## acl_testing5 config
		{
			testName: "acl_segments_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_segments/config_object.json -s %s/acl_segments/conn_spec.json -o %s/acl_segments_tf/nacl_expected.tf",
		},

		// acl testing 5 (json, json single, tf, tf single)
		{
			testName: "acl_testing5_json",
			command:  "../bin/vpcgen synth acl -c %s/acl_testing5/config_object.json -s %s/acl_testing5/conn_spec.json -o %s/acl_testing5_json/nacl_expected.json",
		},
		{
			testName: "acl_testing5_json_single",
			command:  "../bin/vpcgen synth acl --single -c %s/acl_testing5/config_object.json -s %s/acl_testing5/conn_spec.json -o %s/acl_testing5_json_single/nacl_single_expected.json",
		},
		{
			testName: "acl_testing5_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_testing5/config_object.json -s %s/acl_testing5/conn_spec.json -o %s/acl_testing5_tf/nacl_expected.tf",
		},
		{
			testName: "acl_testing5_tf_single",
			command:  "../bin/vpcgen synth acl --single -c %s/acl_testing5/config_object.json -s %s/acl_testing5/conn_spec.json -o %s/acl_testing5_tf_single/nacl_single_expected.tf",
		},

		// acl tg multiple (json, tf, tf separate)
		{
			testName: "acl_tg_multiple_json",
			command:  "../bin/vpcgen synth acl -c %s/acl_tg_multiple/config_object.json -s %s/acl_tg_multiple/conn_spec.json -o %s/acl_tg_multiple_json/nacl_expected.json",
		},
		{
			testName: "acl_tg_multiple_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_tg_multiple/config_object.json -s %s/acl_tg_multiple/conn_spec.json -o %s/acl_tg_multiple_tf/nacl_expected.tf",
		},
		{
			testName: "acl_tg_multiple_tf_separate",
			command:  "../bin/vpcgen synth acl -c %s/acl_tg_multiple/config_object.json -s %s/acl_tg_multiple/conn_spec.json -d %s/acl_tg_multiple_tf_separate -f tf",
		},
	}
}

//nolint:lll // commands can be long
func synthSGTestsList() []mainTestCase {
	return []mainTestCase{
		// sg protocols (all output fmts, externals, scoping, nif as a resource)    ## tg-multiple config
		{
			testName: "sg_protocols_csv",
			command:  "../bin/vpcgen synth sg -c %s/sg_protocols/config_object.json -s %s/sg_protocols/conn_spec.json -o %s/sg_protocols_csv/sg_expected.csv",
		},
		{
			testName: "sg_protocols_json",
			command:  "../bin/vpcgen synth sg -c %s/sg_protocols/config_object.json -s %s/sg_protocols/conn_spec.json -o %s/sg_protocols_json/sg_expected.json",
		},
		{
			testName: "sg_protocols_md",
			command:  "../bin/vpcgen synth sg -c %s/sg_protocols/config_object.json -s %s/sg_protocols/conn_spec.json -o %s/sg_protocols_md/sg_expected.md",
		},
		{
			testName: "sg_protocols_tf",
			command:  "../bin/vpcgen synth sg -c %s/sg_protocols/config_object.json -s %s/sg_protocols/conn_spec.json -o %s/sg_protocols_tf/sg_expected.tf",
		},

		// sg testing 3 (all fmts, VPEs are included)
		{
			testName: "sg_testing3_csv",
			command:  "../bin/vpcgen synth sg -c %s/sg_testing3/config_object.json -s %s/sg_testing3/conn_spec.json -o %s/sg_testing3_csv/sg_expected.csv",
		},
		{
			testName: "sg_testing3_json",
			command:  "../bin/vpcgen synth sg -c %s/sg_testing3/config_object.json -s %s/sg_testing3/conn_spec.json -o %s/sg_testing3_json/sg_expected.json",
		},
		{
			testName: "sg_testing3_md",
			command:  "../bin/vpcgen synth sg -c %s/sg_testing3/config_object.json -s %s/sg_testing3/conn_spec.json -o %s/sg_testing3_md/sg_expected.md",
		},
		{
			testName: "sg_testing3_tf",
			command:  "../bin/vpcgen synth sg -c %s/sg_testing3/config_object.json -s %s/sg_testing3/conn_spec.json -o %s/sg_testing3_tf/sg_expected.tf",
		},

		// sg tg multiple (tf separate)
		{
			testName: "sg_tg_multiple_tf_separate",
			command:  "../bin/vpcgen synth sg -c %s/sg_tg_multiple/config_object.json -s %s/sg_tg_multiple/conn_spec.json -d %s/sg_tg_multiple_tf_separate -f tf",
		},
	}
}
