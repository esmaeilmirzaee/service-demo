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
		-t service-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ------------------------------------------------------------
# Running from within k8s/kind

KIND_CLUSTER := thebeaver-starter-cluster

kind-up:
	kind create cluster \
		--name kindest/node@sha256:32b8b755dee8d5fc6dbafef80898e7f2e450655f45f4b59cca3f7f57e9278c3b \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	# default namespace
	kubectl config set-context --current --namespace=service-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	kind load docker-image service-amd64:$(version) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/base/service-pod/base-service.yaml | kubectl apply
	-f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch

kind-status-service:
	kubectl get pods -o wide --watch

kind-logs:
	kubectl logs -l app=service --all-containers=true -f --tail=100


kind-restart:
	kubectl rollout restart deployment service-pod

kind-update: all kind-load kind-restart


kind-describe:
	kubectl describe pod -l app=service
