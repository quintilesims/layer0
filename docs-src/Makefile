L0_VERSION:=$(shell git describe --tags)
export L0_VERSION

deps:
	pip install mkdocs
	pip install -Iv mkdocs-material==3.0.6
	go get github.com/blang/expenv
	go install github.com/blang/expenv

release: template build

template:
	$(shell sed -E -i '' 's/(v[0-9]+.[0-9]+.[0-9]+)/'$(L0_VERSION)'/g' ../README.md)
	expenv -f templates/index.md.template > docs/index.md
	cat templates/release-table-header.template > docs/releases.md
	expenv -f templates/release-table-insertion.template > templates/tmp
	cat templates/release-table-body.template >> templates/tmp
	mv templates/tmp templates/release-table-body.template
	cat templates/release-table-body.template >> docs/releases.md

build:
	mkdocs build --clean

serve:
	mkdocs serve -a 0.0.0.0:9090

.PHONY: deps release template build serve
