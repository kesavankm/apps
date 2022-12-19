# Where to push the docker image.
REGISTRY ?= docker.io/kesavank7

# This version-strategy uses git tags to set the version string
VERSION ?= $(shell git describe --tags --always --dirty)

GOOS := linux
GOARCH := amd64
APP_NAME := $(notdir $(CURDIR))

# Used internally.  Users should pass GOOS and/or GOARCH.
OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

# Directories that we need created to build/test.
BUILD_DIRS := bin/$(OS)_$(ARCH)                  \
              .go/bin/$(OS)_$(ARCH) \
              .go/cache             \
              .go/pkg

BASEIMAGE ?= golang:1.17-alpine

IAM ?= $(shell whoami)

BUILD_IMAGE ?= $(REGISTRY)/golang:1.17-alpine

.PHONY: all build test clean push docker

$(BUILD_DIRS):
	@mkdir -p $@

# If you want to build all binaries, see the 'all-build' rule.
# If you want to build all containers, see the 'all-container' rule.
# If you want to build AND push all containers, see the 'all-push' rule.
all: # @HELP builds binaries for one platform ($OS/$ARCH)
all: build

hello:
	@echo hello
	@echo $(APP_NAME)

build: $(BUILD_DIRS)
	@echo
	@echo "building for $(OS)/$(ARCH)"
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
		-v ${PWD}:/app								\
		-w /app 									\
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/pkg:/go/pkg                              \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    -v $(HOME)/.ssh:/root/.ssh                              \
	    -v $(HOME)/.gitconfig:/root/.gitconfig                  \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    $(BUILD_IMAGE)                                          \
	    /bin/sh -c "                                            \
	        ARCH=$(ARCH)                                        \
	        OS=$(OS)                                            \
	        VERSION=$(VERSION)                                  \
	        ./build/build_app.sh                             	\
		"

# Run go fmt against code
fmt:
	go mod tidy
	go fmt ./...

DKR_IMG_NAME := app-$(APP_NAME)
DKR_IMG_TAG ?= latest

clean:
	rm -rf .go bin

help:
	@echo "App name: " ${APP_NAME}
	@echo "App version:" ${VERSION}
	@echo "Registry:" ${REGISTRY}
	@echo "Docker image:" ${DKR_IMG_NAME}
	@echo "Docker tag:" ${DKR_IMG_TAG}

# This is the target definition for all container-dotfiles.
# These are used to track build state in hidden files.
docker: # @HELP pushes the container for one platform ($OS/$ARCH) to the defined registry
docker:
	@echo "building docker"
	@docker build -t $(REGISTRY)/$(DKR_IMG_NAME):$(DKR_IMG_TAG) -f build/Dockerfile .
	@docker images -q $(REGISTRY)/$(DKR_IMG_NAME):$(DKR_IMG_TAG) > $@
	@echo

# This is the target definition for all container-dotfiles.
# These are used to track build state in hidden files.
push: # @HELP pushes the container for one platform ($OS/$ARCH) to the defined registry
push:
	@docker build -t $(REGISTRY)/$(DKR_IMG_NAME):$(DKR_IMG_TAG) -f build/Dockerfile .
	@docker images -q $(REGISTRY)/$(DKR_IMG_NAME):$(DKR_IMG_TAG) > $@
	@echo

	docker push $(REGISTRY)/$(DKR_IMG_NAME):$(DKR_IMG_TAG) 

version: # @HELP outputs the version string
version:
	@echo $(VERSION)

# Example: make shell CMD="-c 'date > datefile'"
shell: # @HELP launches a shell in the containerized build environment
shell: $(BUILD_DIRS)
	@echo "launching a shell in the containerized build environment"
	@docker run                                                 \
	    -ti                                                     \
	    --rm                                                    \
		-v ${PWD}:/app								\
		-w /app 									\
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/pkg:/go/pkg                              \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    -v $(HOME)/.ssh:/root/.ssh                              \
	    -v $(HOME)/.gitconfig:/root/.gitconfig                  \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    $(BUILD_IMAGE)                                          \
	    /bin/sh $(CMD)