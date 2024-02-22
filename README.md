# vpc-network-config-synthesis

Tool for automatic synthesis of VPC network configurations, namely Network ACLs and Security Groups.

* spec_schema.json is the JSON schema for VPC-synthesis

Build:

```commandline
make build
```


Test:

```commandline
make test
```


Run on Linux environment:

```commandline
bin/vpcgen -target=acl -config test/data/acl_testing5/config_object.json test/data/acl_testing5/conn_spec.json

$ bin/vpcgen -target=sg -config test/data/sg_testing2/config_object.json test/data/sg_testing2/conn_spec.json
```


Run on Windows:

```commandline
bin\vpcgen.exe -target=acl -config test\data\acl_testing5\config_object.json test\data\acl_testing5\conn_spec.json

bin\vpcgen.exe -target=sg -config test\data\sg_testing2\config_object.json test\data\sg_testing2\conn_spec.json
```




## Code generation

Install [atombender/go-jsonschema](https://github.com/atombender/go-jsonschema)
(important: **not** [xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema))

```commandline
go get github.com/atombender/go-jsonschema/...
go install github.com/atombender/go-jsonschema@latest
```

Then run

```commandline
make generate
```

The result is written into [pkg/io/jsonio/data_model.go](pkg/io/jsonio/data_model.go).
