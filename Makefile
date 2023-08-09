REPOSITORY := github.com/np-guard/vpc-network-config-synthesis
EXE := vpcgen.exe

./bin/$(EXE): build

.PHONY: mod fmt lint generate build test jd-test

mod: go.mod
	@echo -- $@ --
	go mod tidy
	go mod download

fmt:
	@echo -- $@ --
	dos2unix * .* pkg/*/*.go cmd/*/*.go spec_schema.json examples/*
	goimports -local $(REPOSITORY) -w .

lint:
	@echo -- $@ --
	# to avoid parse errors, use git's diff - in windows, add C:\Program Files\Git\usr\bin\ to PATH
	golangci-lint run --new

precommit: mod fmt lint

pkg/synth/data_model.go: spec_schema.json
	@echo -- generate --
	# Install https://github.com/omissis/go-jsonschema
	gojsonschema spec_schema.json --package synth --struct-name-from-title --tags json --output $@
	goimports -local $(REPOSITORY) -w $@

generate: pkg/synth/data_model.go

build:
	@echo -- $@ --
	CGO_ENABLED=0 go build -o ./bin/$(EXE) ./cmd/vpcgen

test:
	@echo -- $@ --
	go test ./... -v -cover -coverprofile synth.coverprofile

jd-test: ./bin/$(EXE)
	@echo -- $@ --
	# Install https://github.com/josephburnett/jd
	./bin/$(EXE) examples/generic_example.json > tmp.json
	jd examples/generic_example.json tmp.json
