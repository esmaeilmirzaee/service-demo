SHELL := ${SHELL}

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
		-f zarf/docker/dockerfile \
		-t service-demo-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ------------------------------------------------------------
# Running from within k8s/kind

KIND_CLUSTER := ardan-starter-cluster

kind-up:
	kind create cluster \
		--name kindest/node:v1.21
		.1@sha256
		:32b8b755dee8d5fc6dbafef80898e7f2e450655f45f4b59cca3f7f57e9278c3b \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
