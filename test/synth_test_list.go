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
		// acl segments (bidi)
		{
			testName: "acl_segments_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_segments_tf/config_object.json -s %s/acl_segments_tf/conn_spec.json -o %s/acl_segments_tf/nacl_expected.tf",
		},

		// acl protocols (all output fmts, externals are included)
		{
			testName: "acl_protocols_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_protocols_tf/config_object.json -s %s/acl_protocols_tf/conn_spec.json -o %s/acl_protocols_tf/nacl_expected.tf",
		},
		{
			testName: "acl_protocols_csv",
			command:  "../bin/vpcgen synth acl -c %s/acl_protocols_csv/config_object.json -s %s/acl_protocols_csv/conn_spec.json -o %s/acl_protocols_csv/nacl_expected.csv",
		},
		{
			testName: "acl_protocols_md",
			command:  "../bin/vpcgen synth acl -c %s/acl_protocols_md/config_object.json -s %s/acl_protocols_md/conn_spec.json -o %s/acl_protocols_md/nacl_expected.md",
		},
		{
			testName: "acl_protocols_json",
			command:  "../bin/vpcgen synth acl -c %s/acl_protocols_json/config_object.json -s %s/acl_protocols_json/conn_spec.json -o %s/acl_protocols_json/nacl_expected.json",
		},

		// acl nif (scoping)
		{
			testName: "acl_nif_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_nif_tf/config_object.json -s %s/acl_nif_tf/conn_spec.json -o %s/acl_nif_tf/nacl_expected.tf",
		},

		// acl testing 5 (tf, tf single, tf separate, tf single separate, json, json single)
		{
			testName: "acl_testing5_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_testing5_tf/config_object.json -s %s/acl_testing5_tf/conn_spec.json -o %s/acl_testing5_tf/nacl_expected.tf",
		},
		{
			testName: "acl_testing5_tf_single",
			command:  "../bin/vpcgen synth acl --single -c %s/acl_testing5_tf_single/config_object.json -s %s/acl_testing5_tf_single/conn_spec.json -o %s/acl_testing5_tf_single/nacl_single_expected.tf",
		},
		{
			testName: "acl_testing5_tf_separate",
			command:  "../bin/vpcgen synth acl -c %s/acl_testing5_tf_separate/config_object.json -s %s/acl_testing5_tf_separate/conn_spec.json -d %s/acl_testing5_tf_separate -f tf",
		},
		{
			testName: "acl_testing5_tf_single_separate",
			command:  "../bin/vpcgen synth acl --single -c %s/acl_testing5_tf_single_separate/config_object.json -s %s/acl_testing5_tf_single_separate/conn_spec.json -d %s/acl_testing5_tf_single_separate -f tf -p single",
		},
		{
			testName: "acl_testing5_json",
			command:  "../bin/vpcgen synth acl -c %s/acl_testing5_json/config_object.json -s %s/acl_testing5_json/conn_spec.json -o %s/acl_testing5_json/nacl_expected.json",
		},
		{
			testName: "acl_testing5_json_single",
			command:  "../bin/vpcgen synth acl --single -c %s/acl_testing5_json_single/config_object.json -s %s/acl_testing5_json_single/conn_spec.json -o %s/acl_testing5_json_single/nacl_single_expected.json",
		},

		// acl tg multiple (tf -l, tf separate, json)
		{
			testName: "acl_tg_multiple_tf",
			command:  "../bin/vpcgen synth acl -c %s/acl_tg_multiple_tf/config_object.json -s %s/acl_tg_multiple_tf/conn_spec.json -o %s/acl_tg_multiple_tf/nacl_expected.tf -l",
		},
		{
			testName: "acl_tg_multiple_separate",
			command:  "../bin/vpcgen synth acl -c %s/acl_tg_multiple_separate/config_object.json -s %s/acl_tg_multiple_separate/conn_spec.json -d %s/acl_tg_multiple_separate -f tf",
		},
		{
			testName: "acl_tg_multiple_json",
			command:  "../bin/vpcgen synth acl -c %s/acl_tg_multiple_json/config_object.json -s %s/acl_tg_multiple_json/conn_spec.json -o %s/acl_tg_multiple_json/nacl_expected.json",
		},
	}
}

//nolint:lll // commands can be long
func synthSGTestsList() []mainTestCase {
	return []mainTestCase{
		// sg protocols (all output fmts, externals, scoping, nif as a resource)
		{
			testName: "sg_protocols_tf",
			command:  "../bin/vpcgen synth sg -c %s/sg_protocols_tf/config_object.json -s %s/sg_protocols_tf/conn_spec.json -o %s/sg_protocols_tf/sg_expected.tf",
		},
		{
			testName: "sg_protocols_csv",
			command:  "../bin/vpcgen synth sg -c %s/sg_protocols_csv/config_object.json -s %s/sg_protocols_csv/conn_spec.json -o %s/sg_protocols_csv/sg_expected.csv",
		},
		{
			testName: "sg_protocols_json",
			command:  "../bin/vpcgen synth sg -c %s/sg_protocols_json/config_object.json -s %s/sg_protocols_json/conn_spec.json -o %s/sg_protocols_json/sg_expected.json",
		},
		{
			testName: "sg_protocols_md",
			command:  "../bin/vpcgen synth sg -c %s/sg_protocols_md/config_object.json -s %s/sg_protocols_md/conn_spec.json -o %s/sg_protocols_md/sg_expected.md",
		},

		// sg testing 3 (all fmts, VPEs are included)
		{
			testName: "sg_testing3_tf",
			command:  "../bin/vpcgen synth sg -c %s/sg_testing3_tf/config_object.json -s %s/sg_testing3_tf/conn_spec.json -o %s/sg_testing3_tf/sg_expected.tf",
		},
		{
			testName: "sg_testing3_csv",
			command:  "../bin/vpcgen synth sg -c %s/sg_testing3_csv/config_object.json -s %s/sg_testing3_csv/conn_spec.json -o %s/sg_testing3_csv/sg_expected.csv",
		},
		{
			testName: "sg_testing3_json",
			command:  "../bin/vpcgen synth sg -c %s/sg_testing3_json/config_object.json -s %s/sg_testing3_json/conn_spec.json -o %s/sg_testing3_json/sg_expected.json",
		},
		{
			testName: "sg_testing3_md",
			command:  "../bin/vpcgen synth sg -c %s/sg_testing3_md/config_object.json -s %s/sg_testing3_md/conn_spec.json -o %s/sg_testing3_md/sg_expected.json",
		},

		// sg tg multiple (tf separate)
		{
			testName: "sg_g_multiple_tf_separate",
			command:  "../bin/vpcgen synth sg -c %s/sg_g_multiple_tf_separate/config_object.json -s %s/sg_g_multiple_tf_separate/conn_spec.json -d %s/sg_g_multiple_tf_separate -f tf",
		},
	}
}
