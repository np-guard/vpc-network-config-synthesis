# vpc-network-config-synthesis
Tool for automatic synthesis of VPC network resources.

* spec_schema.json is the JSON schema for VPC-synthesis

Build:
```commandline
make build
```

Test:
```commandline
make test
```
Run:
```commandline
bin/vpcgen examples/generic-example.json
```

## Code generation

Install [omissis/go-jsonschema](https://github.com/omissis/go-jsonschema) (important: **not** [xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema))
```commandline
go install github.com/omissis/go-jsonschema
```

Then run
```commandline
make generate
```
The result is written into [synth/data_model.go](synth/data_model.go).
