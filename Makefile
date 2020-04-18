
GO111          ?= on
CGO_ENABLED    ?= 0
REPOSITORY     ?= kj187/http-monitoring
REF            ?= master

default: info

info:
	@go version

setup: info ## Initial setup
	GO111MODULE=$(GO111) go mod vendor

lint: info ## Lint the files
	./scripts/go_lint.sh

go_build: info ## Build the binary file
	GO111MODULE=$(GO111) go mod vendor
	GO111MODULE=$(GO111) CGO_ENABLED=$(CGO_ENABLED) go build -o http-monitoring -mod vendor src/cmd/http-monitoring/main.go

docker_build: ## Build docker image
	docker build . --file Dockerfile --tag $(REPOSITORY)

dockerhub_login: ## Login to Dockerhub
	echo "$${DOCKER_REGISTRY_PASSWORD}" | docker login -u $${DOCKER_REGISTRY_ACTOR} --password-stdin

docker_push: ## Push docker image to dockerhub
	REPOSITORY=$(REPOSITORY) REF=$(REF) ./scripts/docker_push.sh

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: help