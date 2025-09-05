# Disable all the default make stuff
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

## Display a list of the documented make targets
.PHONY: help
help:
	@echo Documented Make targets:
	@perl -e 'undef $$/; while (<>) { while ($$_ =~ /## (.*?)(?:\n# .*)*\n.PHONY:\s+(\S+).*/mg) { printf "\033[36m%-30s\033[0m %s\n", $$2, $$1 } }' $(MAKEFILE_LIST) | sort

.PHONY: .FORCE
.FORCE:

build:
	go build ./cmd/score-aca/

test:
	go vet ./...
	go test ./... -cover -race

test-app: build
	./score-aca --version
	./score-aca init
	cat score.yaml
	./score-aca generate score.yaml
	cat manifest.bicep

build-container:
	docker build -t score-aca:local .

test-container: build-container
	docker run --rm score-aca:local --version
	docker run --rm -v .:/score-aca score-aca:local init
	cat score.yaml
	docker run --rm -v .:/score-aca score-aca:local generate score.yaml
	cat manifest.bicep