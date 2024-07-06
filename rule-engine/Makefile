# Rule Engine Make File
# Made available under the Apache-2.0 License

# User Settable Values
OS=$(shell uname -s | tr [:upper:] [:lower:])
ARCH=$(shell uname -m)
GOARCH=$(subst x86_64,amd64,$(ARCH))
EXT=$(if $(filter win%,$(OS)),.exe,)
TESTFILES=./...

GO=go
GO_COMPILE_FLAGS=
SCHEMA_SRC_GEN=go-jsonschema

SOURCES := $(shell find . -name '*.go' -type f)
OUTDIR := out
NAME := qazaar-rule

DOCUMENT_SCHEMA := ../data-exchange/schema/document-description.v1.schema.json
DOCUMENT_PARSER_SRC := schema/document/schema.go
ONTOLOGY_SCHEMA := ../data-exchange/schema/ontology.v1.schema.json
ONTOLOGY_PARSER_SRC := schema/ontology/schema.go
RULES_SCHEMA := ../data-exchange/schema/rules.v1.schema.json
RULES_PARSER_SRC := schema/rules/schema.go

## Run the primary build tasks.
main: build test

## Compile the commands.
build: $(OUTDIR)/$(NAME)-$(OS)-$(ARCH)$(EXT)

## Remove built files.  Does not remove generated code from the schema.
clean: .FORCE
	rm -r $(OUTDIR) || true

## Run tests.
test: .FORCE
	@mkdir -p $(OUTDIR)/profiling
	$(GO) test -outputdir $(OUTDIR)/profiling -coverprofile cover.out $(TESTFILES)

$(OUTDIR)/$(NAME)-$(OS)-$(ARCH)$(EXT): $(OUTDIR) $(SOURCES)
	GOARCH=$(GOARCH) GOOS=$(OS) $(GO) build -o $@ ./cmd

$(OUTDIR):
	mkdir -p $@

## Re-generate the schema parsing sources.
##   This requires that you have run the `dependencies` target and have `$HOME/go/bin` in your path.
schema: $(DOCUMENT_PARSER_SRC) $(ONTOLOGY_PARSER_SRC) $(RULES_PARSER_SRC)

$(DOCUMENT_PARSER_SRC): $(DOCUMENT_SCHEMA)
	@mkdir -p "`dirname "$@"`"
	$(SCHEMA_SRC_GEN) -p document $< > $@ || ( rm $@ ; exit 1 )

$(ONTOLOGY_PARSER_SRC): $(ONTOLOGY_SCHEMA)
	@mkdir -p "`dirname "$@"`"
	$(SCHEMA_SRC_GEN) -p ontology $< > $@ || ( rm $@ ; exit 1 )

$(RULES_PARSER_SRC): $(RULES_SCHEMA)
	@mkdir -p "`dirname "$@"`"
	$(SCHEMA_SRC_GEN) -p rules $< > $@ || ( rm $@ ; exit 1 )

## Install dependencies.
dependencies: .FORCE
	go install github.com/atombender/go-jsonschema@v0.16.0

# General Makefile (not GNU-specific) rule to require a rule to run.
.FORCE:
