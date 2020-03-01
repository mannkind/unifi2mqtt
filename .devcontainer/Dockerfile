# Base Build Image
FROM golang:1.14-alpine3.11
WORKDIR /app
RUN apk add --no-cache --update build-base git mosquitto openssh && \
    git clone https://github.com/magefile/mage && \
    cd mage && go run bootstrap.go install && cd - && \
    rm -rf mage