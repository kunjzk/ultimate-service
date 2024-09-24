SHELL := /bin/bash

run:
	go run main.go

build:
	go build -ldflags "-X main.build=local"

VERSION := 1.0

all: service

service:
	docker build \
		--no-cache \
		-f zarf/docker/dockerfile \
		-t ultimate-service-amd64:$(VERSION)\
		--build-arg BUILD_REF=$(VERSION)\
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"`\
		.

KIND_CLUSTER := ardan-starter-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yml
	kubectl config set-context --current --namespace=service-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	kind load docker-image ultimate-service-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	cat zarf/k8s/base/service-pod/base-service.yml | kubectl apply -f -

kind-logs:
	kubectl logs -l app=service --all-containers=true -f --tail=100

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-restart:
	kubectl rollout restart deployment service-pod

kind-status-service:
	kubectl get pods -o wide --watch

kind-update: all kind-load kind-restart

kind-describe:
	kubectl describe pod -l app=service