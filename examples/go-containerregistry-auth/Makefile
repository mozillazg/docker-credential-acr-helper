
IMAGE ?= mozillazg/get-acr-digest:0.1.0

help:
	@echo "build-image              build image"
	@echo "push-image               push image"
	@echo "build-push-image         build and push image"


build-push-image: build-image push-image

push-image:
	docker push $(IMAGE)

build-image:
	docker build -t $(IMAGE) -f Dockerfile ../../
