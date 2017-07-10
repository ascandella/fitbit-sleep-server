IMAGE_NAME := ai-life
CONTAINER_NAME ?= $(IMAGE_NAME)
CREDENTIALS ?= /etc/ai-life

.PHONY: build-docker
build-docker:
	docker build --tag '$(IMAGE_NAME)' .

.PHONY: run-docker
run-docker:
	docker run \
		--name $(CONTAINER_NAME) \
		--volume $(CREDENTIALS):/config \
		--publish 3030:3030 \
		$(IMAGE_NAME)

.PHONY: kill-docker
kill-docker:
	docker stop -t 15 $(CONTAINER_NAME)
	docker rm $(CONTAINER_NAME)

.PHONY: test
test:
	@go test -v -cover .
	@go vet .
	@golint .

.DEFAULT_GOAL := test
