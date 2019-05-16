GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
WIRECMD=$$GOPATH/bin/wire gen
BINARY_NAME=unifi2mqtt
BINARY_VERSION=$(shell git describe --tags --always --dirty="-dev")
BINARY_DATE=$(shell date -u '+%Y-%m-%d-%H%M UTC')
BINARY_VERSION_FLAGS=-ldflags='-X "main.Version=$(BINARY_VERSION)" -X "main.BuildTime=$(BINARY_DATE)"'
DOCKER_IMAGE=mannkind/unifi2mqtt
DOCKER_ARCHS=amd64 arm32v6 arm64v8
DOCKER_VERSION=$(BINARY_VERSION)

all: clean wire build test format vet
test: 
		$(GOTEST) --coverprofile=/tmp/app.cover -v ./...
format:
	    $(GOFMT) .
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		{ \
		for arch in $(DOCKER_ARCHS); do \
			rm -f Dockerfile.$${arch} ;\
		done ;\
		}
vet:
	    $(GOVET) .
get_wire:
		$(GOGET) github.com/google/wire/cmd/wire
wire: get_wire
		$(WIRECMD)
build: format
		$(GOBUILD) $(BINARY_VERSION_FLAGS) -o $(BINARY_NAME) -v
run: build
		./$(BINARY_NAME)
docker: clean
		{ \
		set -e ;\
		for arch in $(DOCKER_ARCHS); do \
		case $${arch} in \
			amd64   ) golang_arch="amd64";; \
			arm32v6 ) golang_arch="arm";; \
			arm64v8 ) golang_arch="arm64";; \
		esac ;\
		cp Dockerfile.template Dockerfile.$${arch} ;\
		sed -i "" "s|__BASEIMAGE_ARCH__|$${arch}|g" Dockerfile.$${arch} ;\
		sed -i "" "s|__GOLANG_ARCH__|$${golang_arch}|g" Dockerfile.$${arch} ;\
		done ;\
		}

		$(foreach arch,$(DOCKER_ARCHS),docker build --no-cache --pull -q -f Dockerfile.$(arch) -t $(DOCKER_IMAGE):$(arch)-$(DOCKER_VERSION) . ;)
docker-push:
		$(foreach arch,$(DOCKER_ARCHS),docker push $(DOCKER_IMAGE):$(arch)-$(DOCKER_VERSION);)
		docker manifest create $(DOCKER_IMAGE):$(DOCKER_VERSION) $(DOCKER_IMAGE):amd64-$(DOCKER_VERSION) $(DOCKER_IMAGE):arm32v6-$(DOCKER_VERSION) $(DOCKER_IMAGE):arm64v8-$(DOCKER_VERSION)
		docker manifest annotate $(DOCKER_IMAGE):$(DOCKER_VERSION) $(DOCKER_IMAGE):arm32v6-$(DOCKER_VERSION) --os linux --arch arm --variant v6
		docker manifest annotate $(DOCKER_IMAGE):$(DOCKER_VERSION) $(DOCKER_IMAGE):arm64v8-$(DOCKER_VERSION) --os linux --arch arm64 --variant v8
		docker manifest push --purge $(DOCKER_IMAGE):$(DOCKER_VERSION)
