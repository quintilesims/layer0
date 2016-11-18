deps:
	pip install awscli

update:
	aws s3 sync s3://docs.xfra.ims.io ./site

dev: deps
	pip install mkdocs

build:
	mkdocs build --clean
	
deploy: build
	aws s3 sync site s3://docs.xfra.ims.io --delete
