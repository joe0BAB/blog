GO_VERSION := $(shell sh -c "awk '/^go / { print \$$2 }' go.mod")

ifeq ($(OS),Windows_NT)
	EXTENSION = .exe
else
	EXTENSION =
endif

# golangci-lint must be pinned - linters can become more strict on upgrade
GOLANGCI_LINT_VERSION := v1.55.2
GO_LDFLAGS =
export GO_VERSION GO_LDFLAGS GOPRIVATE GOLANGCI_LINT_VERSION

BUILDER=buildx-multiarch

DOCKER_BUILD_ARGS := --ssh default \
					--build-arg GO_VERSION \
					--build-arg GO_LDFLAGS \
          			--build-arg GOLANGCI_LINT_VERSION

LINT_PLATFORMS = linux,darwin,windows

multiarch-builder: ## Create buildx builder for multi-arch build, if not exists
	docker buildx inspect $(BUILDER) || docker buildx create --name=$(BUILDER) --driver=docker-container --driver-opt=network=host

format: ## Format code
	@docker buildx build $(DOCKER_BUILD_ARGS) -o . --target=format .

lint: multiarch-builder ## Lint code
	@docker buildx build $(DOCKER_BUILD_ARGS) --pull --builder=$(BUILDER) --target=check-mod .
	@docker buildx build $(DOCKER_BUILD_ARGS) --pull --builder=$(BUILDER) --target=lint --platform=$(LINT_PLATFORMS) .

clean: ## remove built binaries and packages
	@sh -c "rm -rf bin dist"

build:
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags "-s -w ${GO_LDFLAGS}" -o ./bin/blog$(EXTENSION) ./cmd
