REPOSITORY := github.com/np-guard/vpc-network-config-synthesis
ifeq ($(OS),Windows_NT)
    SETCGO = set
	TARGETNAME = vpcgen.exe
else
	SETCGO =
    TARGETNAME = vpcgen
endif
TARGET = ./bin/$(TARGETNAME)

$(TARGET): build

.PHONY: mod fmt lint generate build test jd-test

mod: go.mod
	@echo -- $@ --
	go mod tidy
	go mod download

fmt:
	@echo -- $@ --
	dos2unix * .* pkg/*/*.go cmd/*/*.go spec_schema.json examples/* test/data/*/*
	goimports -local $(REPOSITORY) -w .
	terraform fmt -recursive

lint-go:
	@echo -- $@ --
	# to avoid parse errors, use git's diff - in windows, add C:\Program Files\Git\usr\bin\ to PATH
	golangci-lint run

lint-json:
	@echo -- $@ --
	check-jsonschema test/data/*/conn_spec.json examples/generic_example.json --schemafile spec_schema.json

lint: lint-go lint-json

precommit: mod fmt lint

pkg/io/jsonio/data_model.go: spec_schema.json
	@echo -- generate --
	# Install https://github.com/atombender/go-jsonschema
	go-jsonschema spec_schema.json --package jsonio --struct-name-from-title --tags json --output $@
	goimports -local $(REPOSITORY) -w $@

generate: pkg/io/jsonio/data_model.go

build:
	@echo -- $@ --
	$(SETCGO) CGO_ENABLED=0 go build -o $(TARGET) ./cmd/vpcgen

test:
	@echo -- $@ --
	go test ./... -v -cover -coverprofile synth.coverprofile

jd-test: $(TARGET)
	@echo -- $@ --
	# Install https://github.com/josephburnett/jd
	$(TARGET) examples/generic_example.json > tmp.json
	jd examples/generic_example.json tmp.json
