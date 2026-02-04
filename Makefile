.PHONY: help fmt lint build run test ci clean docker-build docker-push

APP_NAME ?= biya-exporter
BIN_DIR ?= bin
BIN ?= $(BIN_DIR)/$(APP_NAME)

# docker
IMAGE ?= biya/biya-exporter
TAG ?= dev

GO ?= go
GOFLAGS ?=

help:
	@echo "Targets:"
	@echo "  fmt           - gofmt all go files"
	@echo "  lint          - go vet 静态检查"
	@echo "  build         - build binary into ./bin"
	@echo "  run           - run exporter (CONFIG=... optional)"
	@echo "  test          - run unit tests"
	@echo "  ci            - 完整 CI：fmt + lint + test + build 验证"
	@echo "  clean         - remove build artifacts"
	@echo "  docker-build  - build docker image"
	@echo "  docker-push   - push docker image (requires docker login)"

fmt:
	@echo "==> gofmt"
	@gofmt -w cmd internal

lint:
	@echo "==> go vet"
	@$(GO) vet ./...

ci:
	@./scripts/ci.sh

build:
	@echo "==> build $(BIN)"
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 $(GO) build $(GOFLAGS) -trimpath -ldflags "-s -w \
		-X main.version=$(TAG) \
		-X main.commit=$$(git rev-parse --short HEAD 2>/dev/null || echo none)" \
		-o $(BIN) ./cmd/exporter

run:
	@echo "==> run"
	@$(GO) run $(GOFLAGS) ./cmd/exporter -config $(CONFIG)

test:
	@echo "==> test"
	@$(GO) test $(GOFLAGS) ./...

clean:
	@echo "==> clean"
	@rm -rf $(BIN_DIR)

docker-build:
	@echo "==> docker build $(IMAGE):$(TAG)"
	@docker build -t $(IMAGE):$(TAG) .

docker-push:
	@echo "==> docker push $(IMAGE):$(TAG)"
	@docker push $(IMAGE):$(TAG)


