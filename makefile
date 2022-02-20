SHELL := /bin/zsh

run:
	go run main.go

build:
	go build -ldflags "-X main.build=local"

# ++++++++++++++++++++++++++++++++++++++++++++++++++++
# Building container

# $(shell git rev-parse --short HEAD)
VERSION := 1.0

all: service

service:
	docker build \
		-f zarf/dockerfile \
		-t service-demo-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
