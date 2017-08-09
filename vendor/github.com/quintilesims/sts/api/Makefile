DOCKER_IMAGE = "quintilesims/sts"

release:
	docker build -t $(DOCKER_IMAGE) .
	docker push $(DOCKER_IMAGE)

.PHONY: release
