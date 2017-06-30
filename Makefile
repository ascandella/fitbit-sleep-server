.PHONY: build-docker
build-docker:
	docker build --label ai-life --tag 'ai-life' .

.PHONY: run-docker
run-docker:
	docker run --name ai-life \
		--publish 3030:3030 \
		ai-life

.PHONY: kill-docker
kill-docker:
	docker stop -t 15 ai-life
	docker rm ai-life
