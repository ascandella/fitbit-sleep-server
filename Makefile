IMAGE_NAME := ai-life
CONTAINER_NAME ?= $(IMAGE_NAME)

.PHONY: build-docker
build-docker:
	docker build --tag '$(IMAGE_NAME)' .

.PHONY: run-docker
run-docker:
	docker run \
		--name $(CONTAINER_NAME) \
		--publish 3030:3030 \
		$(IMAGE_NAME)

.PHONY: kill-docker
kill-docker:
	docker stop -t 15 $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

.PHONY: test
test:
	@go test -v .

.DEFAULT_GOAL := test
