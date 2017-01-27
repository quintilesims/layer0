deps:
	pip install awscli mkdocs

build:
	mkdocs build --clean

deploy: build
	aws --region us-west-2 s3 sync site s3://docs.xfra.ims.io --delete

gh-deploy:
	mkdocs gh-deploy

.PHONY: deps build deploy gh-deploy
