# vpc-network-config-synthesis

## About the vpc-network-config-synthesis
Tool for automatic synthesis of VPC network configurations, namely Network ACLs and Security Groups.
The support is in multi-vpc input, but all connectivity must be within the same vpc.

## Usage

### nACLs Generation
There is an option to generate an nACL for each subnet separately, or to generate a single nACL for all subnets in the same VPC.
The input supports subnets, subnet segments and CIDR segments.
Note: The segments are defined in the `conn_spec.json` file.

### SGs Generation
The input supports Instances (VSIs), NIFs and VPEs.
Note: If we have created a SG for a VSI (or its specific NIF), it will be applied to all the NIFs of the VSI. The same goes for ReservedIPs of VPE.

### Output
1. If a path is given to the `output-dir` flag, a folder will be created. There we will create several files (as the number of vpcs), where each file contains the collection relevant to the vpc. The file names will be `prefix_vpc`, where prefix is ​​the value received in the `prefix` flag. If no value is passed in the `prefix` flag, the file names will be the vpc names.
2. If path is given to the `output-file` flag, the collection will be written to the given path.
3. If path is neither given to the `output-file` flag nor the `output-dir` flag, the collection will be written to stdout.
Note: in options 1 and 3 it is mandatory to bring a value to the `fmt` flag.

### Global flags
```commandline
VpcGen translates connectivity spec to network ACLs or Security Groups.
Usage:
        vpcgen [flags] SPEC_FILE

SPEC_FILE: JSON file containing connectivity spec, and segments.

Flags:
  -config string
        JSON file containing config spec
  -fmt string
        Output format. One of "tf", "csv", "md"; must not contradict output file suffix. (default "acl").
  -output-dir string
        Output Directory. If unspecified, output will be written to one file.
  -output-file string
        Output to file. If specified, also determines output format.
  -prefix string
        The prefix of the files that will be created.
  -target string
        Target resource to generate. One of "acl", "sg", "singleacl". (default "acl")
```

## Build the project
Make sure you have golang 1.22+ on your platform.

```commandline
git clone git@github.com:np-guard/vpc-network-config-synthesis.git
cd vpc vpc-network-config-synthesis
make mod
make build
```

Note: Windows environment users will run `make build-windows` instead of `make build`


## Run an example
Linux environment:

```commandline
bin/vpcgen -target=acl -config test/data/acl_testing5/config_object.json test/data/acl_testing5/conn_spec.json

bin/vpcgen -target=sg -config test/data/sg_testing3/config_object.json test/data/sg_testing3/conn_spec.json
```

Note: Windows environment users will replace all `/` with `\`