SHELL := /bin/bash

# ====================================================================================
# Testing running system

# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# hey -m GET -c 100 -n 10000 http://localhost:3000/v1/test
# ====================================================================================

# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa keygen_ bits: 2048
# openssl rsa -pubout -in private.pem -out public.pem
# . /sales-admin genkey

run:
	go run app/services/sales-api/main.go | ./r_test

admin:
	go run app/tooling/admin/main.go

# ====================================================================================
#Building containers

VERSION=1.1

all: sales-arm

sales:
	@docker buildx build \
		--platform linux/amd64 \
		-f zarf/docker/Dockerfile \
		-t sales-api-amd64:${VERSION} \
		--build-arg "BUILD_REF=${VERSION}" \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

sales-arm:
	@docker build \
		-f zarf/docker/Dockerfile.sales-api \
		-t sales-api-arm64:${VERSION} \
		--build-arg "BUILD_REF=${VERSION}" \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.


# ====================================================================================
# Running from withing k8s/kind

KIND_CLUSTER:= ardan-starter-cluster

kind-up:
	@kind create cluster \
    --name=$(KIND_CLUSTER) \
    --config zarf/k8s/kind/kind-config.yaml
	@kubectl config set-cluster kind-$(KIND_CLUSTER)
	@kubectl config set-context --current --namespace=sales-system

kind-down:
	@kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	@cd zarf/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-arm64:$(VERSION)
	@kind load docker-image sales-api-arm64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	@kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -

kind-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go

kind-restart:
	kubectl rollout restart deployment sales-pod

kind-status:
	@kubectl get nodes -o wide
	@echo  
	@echo ==================================================================================================================================
	@kubectl get svc -o wide
	@echo  
	@echo ==================================================================================================================================
	@kubectl get pods -o wide --watch -A

kind-status-sales:
	@kubectl get pods -o wide --watch

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-describe:
	@kubectl describe nodes
	@kubectl describe svc
	@kubectl describe pod -l app=sales

# ====================================================================================
# Modules support

tidy:
	@go mod tidy
	@go mod vendor