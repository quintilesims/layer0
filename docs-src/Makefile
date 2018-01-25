L0_VERSION:=$(shell git describe --tags)
export L0_VERSION

deps:
	pip install mkdocs
	pip install -Iv mkdocs-material==2.3.0
	go get github.com/blang/expenv
	go install github.com/blang/expenv

build:
	mkdocs build --clean

release:
	$(shell sed -E -i '' 's/(v[0-9].[0-9]+.[0-9])/'$(L0_VERSION)'/g' ../README.md)
	expenv -f index.md.template > docs/index.md
	expenv -f release_table.template >> docs/releases.md
	
serve:
	mkdocs serve -a 0.0.0.0:9090

.PHONY: deps build serve
