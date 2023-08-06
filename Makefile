REPOSITORY := github.com/np-guard/vpc-network-config-synthesis
EXE := vpcgen

./bin/$(EXE): build

.PHONY: mod fmt lint generate build test jd-test

mod: go.mod
	go mod tidy
	go mod download

fmt:
	goimports -local $(REPOSITORY) -w .

lint:
	golangci-lint run --new

precommit: mod fmt lint

synth/data_model.go: spec_schema.json
	# Install https://github.com/omissis/go-jsonschema
	gojsonschema spec_schema.json --package synth --struct-name-from-title --tags json --output $@

generate: pkg/synth

build:
	CGO_ENABLED=0 go build -o ./bin/$(EXE) ./cmd/main.go

test:
	go test ./... -v -cover -coverprofile synth.coverprofile

jd-test: ./bin/$(EXE)
	# Install https://github.com/josephburnett/jd
	./bin/$(EXE) examples/generic_example.json > tmp.json
	jd examples/generic_example.json tmp.json
