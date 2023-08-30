# Makefile

# Variables
GATEWAY_IMAGE_NAME = gateway-svc
WORKER_IMAGE_NAME = worker-svc

# Targets
.PHONY: test
test:
	go test ./...

.PHONY: build-gateway
build-gateway:
	go build -o dist/$(GATEWAY_IMAGE_NAME) cmd/gateway/main.go

.PHONY: build-worker
build-worker:
	go build -o dist/$(WORKER_IMAGE_NAME) cmd/worker/main.go

.PHONY: docker-build-gateway
docker-build-gateway:
	docker build -t $(GATEWAY_IMAGE_NAME) -f cmd/gateway/Dockerfile .

.PHONY: docker-build-worker
docker-build-worker:
	docker build -t $(WORKER_IMAGE_NAME) -f cmd/worker/Dockerfile .

.PHONY: publish-gateway
publish-gateway: docker-build-gateway
	docker tag $(GATEWAY_IMAGE_NAME) dockerhub/$(GATEWAY_IMAGE_NAME):latest
	docker push your-docker-registry/$(GATEWAY_IMAGE_NAME):latest

.PHONY: publish-worker
publish-worker: docker-build-worker
	docker tag $(WORKER_IMAGE_NAME) dockerhub/$(WORKER_IMAGE_NAME):latest
	docker push your-docker-registry/$(WORKER_IMAGE_NAME):latest
