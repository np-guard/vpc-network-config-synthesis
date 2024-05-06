# vpc-network-config-synthesis

Tool for automatic synthesis of VPC network configurations, namely Network ACLs and Security Groups.

Build:

```commandline
make build
```

Build on Windows:

```commandline
make build-windows
```


Test:

```commandline
make test
```


Run on Linux environment:

```commandline
bin/vpcgen -target=acl -config test/data/acl_testing5/config_object.json test/data/acl_testing5/conn_spec.json

bin/vpcgen -target=sg -config test/data/sg_testing3/config_object.json test/data/sg_testing3/conn_spec.json
```


Run on Windows:

```commandline
bin\vpcgen.exe -target=acl -config test\data\acl_testing5\config_object.json test\data\acl_testing5\conn_spec.json

bin\vpcgen.exe -target=sg -config test\data\sg_testing3\config_object.json test\data\sg_testing3\conn_spec.json
```
