L0_VERSION:=$(shell git describe --tags)
export L0_VERSION_V:=""

deps:
	pip install mkdocs
	pip install -Iv mkdocs-material==2.3.0
	go get github.com/blang/expenv
	go install github.com/blang/expenv

build:
	mkdocs build --clean

release:
	# Remove v from Version
	$(eval L0_VERSION_V := $(shell echo $(L0_VERSION) | cut -c 2-))

	# Update Version to Latest and clean up
	expenv -f index.md.template > docs/index.md
	expenv -f README.md.template > ../README.md
	expenv -f release_table.template >> docs/releases.md

serve:
	mkdocs serve -a 0.0.0.0:9090

.PHONY: deps build serve
