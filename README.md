# vpc-network-config-synthesis

## About vpc-network-config-synthesis
Tool for automatic synthesis of VPC network configurations, namely Network ACLs and Security Groups.
Multi-vpc input is supported. Required connections cannot cross vpc boundaries.

## Usage
Use the `vpcgen` CLI tool with one of the following commands to specify the type of network resources to generate.
* `vpcgen synth sg` - generate Security Groups.
* `vpcgen synth acl` - generate an nACL for each subnet separately.
* `vpcgen synth acl --single` - generate a single nACL for all subnets in the same VPC.

### nACLs Generation
Specifying the `--single` flag results in generating a single nACL for all subnets in the same VPC. Otherwise, an nACL is generated for each subnet separately.
The input supports subnets, subnet segments, CIDR segments, NIFs, instances (VSIs) and externals.
**Note 1**: The segments are defined in the `conn_spec.json` file.
**Note 2**: A required connection between NIFs or VSIs implies connectivity will be allowed between the subnets they are contained in.

### SGs Generation
The input supports Instances (VSIs), NIFs, VPEs and externals.
**Note**: A Security Group, generated for a specific VSI (or for one of its NIFs), will be applied to all the NIFs of the VSI. The same goes for Reserved IPs of a VPE.

### Output
1. If the `output-dir` flag is used, the specified folder will be created, containing one file per VPC. Each generated file will contain the network resources (Security Groups or Network ACLs) relevant to its VPC. File names are set as `prefix_vpc`, where prefix is ​​the value received in the `prefix` flag. If the `prefix` flag is omitted, file names will match VPC names.
2. If the `output-file` flag is used, all generated resources will be written to the specified file.
3. if both `output-file` and `output-dir` flags are not used, the collection will be written to stdout.

### Global options
```commandline
Flags:
  -c, --config string        JSON file containing a configuration object of existing resources
  -f, --format string        Output format; must be one of [tf, csv, md, json]
  -h, --help                 help for vpc-synthesis
  -l, --locals               whether to generate a locals.tf file (only possible when the output format is tf)
  -d, --output-dir string    Write generated resources to files in the specified directory, one file per VPC.
  -o, --output-file string   Write all generated resources to the specified file.
  -p, --prefix string        The prefix of the files that will be created.
  -s, --spec string          JSON file containing spec file
```
**Note**: The infrastructure configuration must always be provided using the `--config` flag.

## Build the project
Make sure you have golang 1.23+ on your platform.

```commandline
git clone git@github.com:np-guard/vpc-network-config-synthesis.git
cd vpc vpc-network-config-synthesis
make mod
make build
```

**Note**: Windows environment users should run `make build-windows` instead of `make build`.


## Run an example

```commandline
bin/vpcgen synth acl -c test/data/acl_testing5/config_object.json -s test/data/acl_testing5/conn_spec.json

bin/vpcgen synth sg -c test/data/sg_testing3/config_object.json -s test/data/sg_testing3/conn_spec.json
```

**Note**: Windows environment users should replace all `/` with `\`.