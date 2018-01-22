L0_VERSION:=$(shell git describe --tags)
DOWNLOAD_LINK:=| v0.10.4 | \
					[macOS](https://github.com/quintilesims/layer0/releases/download/v0.10.4/macOS.zip) | \
					[Linux](https://github.com/quintilesims/layer0/releases/download/v0.10.4/Linux.zip) | \
					[Windows](https://github.com/quintilesims/layer0/releases/download/v0.10.4/Windows.zip

deps:
	pip install mkdocs
	pip install -Iv mkdocs-material==2.3.0

build:
	mkdocs build --clean

release:
	# Remove v from Version
	$(eval L0_VERSION := $(shell echo $(L0_VERSION) | cut -c 2-))

	# Update Version to Latest and clean up
	$(shell sed -E -i '' 's/([0-9].[0-9]+.[0-9])/'$(L0_VERSION)'/g' ../README.md docs/index.md)
	
	# Add new version to release
	$(shell ex -sc '3i|$(DOWNLOAD_LINK)' -cx docs/releases.md)
	$(shell sed -E -i '' '3 s/([0-9].[0-9]+.[0-9])/'$(L0_VERSION)'/' docs/releases.md)

serve:
	mkdocs serve -a 0.0.0.0:9090

.PHONY: deps build serve
