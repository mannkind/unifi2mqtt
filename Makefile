BINARY_BASE_VERSION=v0.8
BINARY_NAME=unifi2mqtt
DOCKER_IMAGE=mannkind/$(BINARY_NAME)

BINARY_VERSION:=$(shell date +$(BINARY_BASE_VERSION).%y%j.%H%M)
DOCKER_VERSION_ADDTL?=
BINARY_VERSION_FLAGS=-ldflags='-X "main.Version=$(BINARY_VERSION)"'
DOCKER_VERSION?=$(BINARY_VERSION)$(DOCKER_VERSION_ADDTL)
DOCKER_LATEST= latest
ifdef DOCKER_VERSION_ADDTL
	DOCKER_LATEST=
endif

all: wire format vet build test
clean: 
	go clean
	rm -f $(BINARY_NAME)
	rm -f Dockerfile.a*
format:
	go fmt .
vet:
	go vet .
get_wire:
	go get github.com/google/wire/cmd/wire
wire: get_wire
	wire gen
build: 
	go build $(BINARY_VERSION_FLAGS) -o $(BINARY_NAME) -v
test: build
	go test --coverprofile=/tmp/app.cover -v ./...
run: build
	./$(BINARY_NAME)
tag:
	git tag -f $(BINARY_VERSION)
tag-push: tag
	sed -i -e "s/github.com/mannkind:$$GITHUB_TOKEN@github.com/g" .git/config
	git push --tags
docker: build clean
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
docker-push: docker
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
