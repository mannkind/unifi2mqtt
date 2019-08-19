BINARY_BASE_VERSION=0.6
BINARY_NAME=unifi2mqtt
DOCKER_IMAGE=mannkind/$(BINARY_NAME)

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
WIRECMD=wire gen
BINARY_VERSION:=$(shell date +$(BINARY_BASE_VERSION).%y%j.%H%M)
DOCKER_VERSION_ADDTL?=
BINARY_VERSION_FLAGS=-ldflags='-X "main.Version=$(BINARY_VERSION)"'
DOCKER_VERSION?=$(BINARY_VERSION)$(DOCKER_VERSION_ADDTL)
DOCKER_LATEST= latest
ifdef DOCKER_VERSION_ADDTL
	DOCKER_LATEST=
endif

all: clean wire build test format vet
test: 
		$(GOTEST) --coverprofile=/tmp/app.cover -v ./...
format:
	    $(GOFMT) .
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
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
	for arch in amd64 arm32v6 arm64v8; do \
		case $${arch} in \
			amd64   ) golang_arch="amd64";; \
			arm32v6 ) golang_arch="arm";; \
			arm64v8 ) golang_arch="arm64";; \
		esac ;\
	  cp Dockerfile.template Dockerfile.$${arch} && \
	  sed -i"" -e "s|__BASEIMAGE_ARCH__|$${arch}|g" Dockerfile.$${arch} && \
	  sed -i"" -e "s|__GOLANG_ARCH__|$${golang_arch}|g" Dockerfile.$${arch} && \
	  docker build --pull -f Dockerfile.$${arch} -t $(DOCKER_IMAGE):$${arch}-$(DOCKER_VERSION) . && \
	  docker tag $(DOCKER_IMAGE):$${arch}-${DOCKER_VERSION} $(DOCKER_IMAGE):$${arch}-latest && \
	  rm -f Dockerfile.$${arch}* ;\
	done

docker-push:
	for VERSION in $(DOCKER_VERSION) $(DOCKER_LATEST); do \
		docker push $(DOCKER_IMAGE):amd64-$${VERSION} && \
		docker push $(DOCKER_IMAGE):arm32v6-$${VERSION} && \
		docker push $(DOCKER_IMAGE):arm64v8-$${VERSION} && \
		docker manifest create $(DOCKER_IMAGE):$${VERSION} \
			$(DOCKER_IMAGE):amd64-$${VERSION} \
			$(DOCKER_IMAGE):arm32v6-$${VERSION} \
			$(DOCKER_IMAGE):arm64v8-$${VERSION} && \
		docker manifest annotate $(DOCKER_IMAGE):$${VERSION} \
			$(DOCKER_IMAGE):arm32v6-$${VERSION} --os linux --arch arm --variant v6 && \
		docker manifest annotate $(DOCKER_IMAGE):$${VERSION} \
			$(DOCKER_IMAGE):arm64v8-$${VERSION} --os linux --arch arm64 --variant v8 && \
		docker manifest push --purge $(DOCKER_IMAGE):$${VERSION} ;\
	done
