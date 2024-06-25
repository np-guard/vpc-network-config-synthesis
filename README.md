# vpc-network-config-synthesis

## About vpc-network-config-synthesis
Tool for automatic synthesis of VPC network configurations, namely Network ACLs and Security Groups.
Multi-vpc input is supported. Required connections cannot cross vpc boundaries.

## Usage
Use the `vpcgen` CLI tool with one of the following commands to specify the type of network resources to generate.
* `vpcgen sg` - generate Security Groups.
* `vpcgen acl` - generate an nACL for each subnet separately.
* `vpcgen acl --single` - generate a single nACL for all subnets in the same VPC.

### nACLs Generation
There is an option to generate an nACL for each subnet separately, or to generate a single nACL for all subnets in the same VPC.
The input supports subnets, subnet segments, CIDR segments and externals.
Note: The segments are defined in the `conn_spec.json` file.

### SGs Generation
The input supports Instances (VSIs), NIFs, VPEs and externals.
**Note**: A Security Group, generated for a specific VSI (or for one of its NIFs), will be applied to all the NIFs of the VSI. The same goes for Reserved IPs of a VPE.

### Output
1. If the `output-dir` flag is used, the specified folder will be created, containing one file per VPC. Each generated file will contain the network resources (Security Groups or Network ACLs) relevant to its VPC. File names are set as `prefix_vpc`, where prefix is ​​the value received in the `prefix` flag. If the `prefix` flag is omitted, file names will match VPC names.
2. If the `output-file` flag is used, all generated resources will be written to the specified file.
3. if both `output-file` and `output-dir` flags are not used, the collection will be written to stdout.

### Global options
```commandline
Tool for automatic synthesis of VPC network configurations, namely Network ACLs and Security Groups.

Usage:
  vpc-synthesis [command]

Available Commands:
  acl         nACL generation for subnets
  sg          SG generation for nifs and vpes

Flags:
  -c, --config string        JSON file containing config spec
  -f, --format string        Output format; must be one of [tf, csv, md, json]
  -h, --help                 help for vpc-synthesis
  -d, --output-dir string    Write generated resources to files in the specified directory, one file per VPC.
  -o, --output-file string   Write all generated resources to the specified file.
      --prefix string        The prefix of the files that will be created.
  -s, --spec string          JSON file containing spec file

```

## Build the project
Make sure you have golang 1.22+ on your platform.

```commandline
git clone git@github.com:np-guard/vpc-network-config-synthesis.git
cd vpc vpc-network-config-synthesis
make mod
make build
```

**Note**: Windows environment users should run `make build-windows` instead of `make build`.


## Run an example

```commandline
bin/vpcgen acl -c test/data/acl_testing5/config_object.json -s test/data/acl_testing5/conn_spec.json

bin/vpcgen sg -c test/data/sg_testing3/config_object.json -s test/data/sg_testing3/conn_spec.json
```

**Note**: Windows environment users should replace all `/` with `\`.