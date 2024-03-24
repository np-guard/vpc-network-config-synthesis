REPOSITORY := github.com/np-guard/vpc-network-config-synthesis
ifeq ($(OS),Windows_NT)
	SET = set
	TARGETNAME = vpcgen.exe
else
	SET =
	TARGETNAME = vpcgen
endif
TARGET = ./bin/$(TARGETNAME)

$(TARGET): build

.PHONY: mod fmt lint build test

mod: go.mod
	@echo -- $@ --
	go mod tidy
	go mod download

fmt:
	@echo -- $@ --
	dos2unix * .* pkg/*/*.go cmd/*/*.go examples/* test/data/*/*
	goimports -local $(REPOSITORY) -w .
	terraform fmt -recursive

lint-go:
	@echo -- $@ --
	# to avoid parse errors, use git's diff - in windows, add C:\Program Files\Git\usr\bin\ to PATH
	golangci-lint run

lint: lint-go

precommit: mod fmt lint

build:
	@echo -- $@ --
	$(SET) CGO_ENABLED=0 go build -o $(TARGET) ./cmd/vpcgen

test:
	@echo -- $@ --
	go test ./... -v -cover -coverprofile synth.coverprofile
