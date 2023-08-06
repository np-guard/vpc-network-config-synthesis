REPOSITORY := github.com/np-guard/vpc-network-config-synthesis
EXE := vpcsynthesizer

mod: go.mod
	go mod tidy
	go mod download

fmt:
	goimports -local $(REPOSITORY) -w .

lint:
	golangci-lint run --new

precommit: mod fmt lint

synth/data_model.go: spec_schema.json
	gojsonschema spec_schema.json --package synth --struct-name-from-title --tags json --output $@

generate: synth/data_model.go

build:
	CGO_ENABLED=0 go build -o ./bin/$(EXE) ./cmd/main.go

./bin/$(EXE): build

test: ./bin/$(EXE)
	./bin/$(EXE) examples/generic_example.json > tmp.json
	jd examples/generic_example.json tmp.json
	#go test ./... -v -cover -coverprofile analyzer.coverprofile
