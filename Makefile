.PHONY: build

BIN_DIR := bin
UIKIT_BUNDLE := $(BIN_DIR)/edge-spark-ui-kit-src.tar.gz

build:
	@mkdir -p $(BIN_DIR)
	@go test ./...
	@tar -czf $(UIKIT_BUNDLE) README.md go.mod go.sum uikit.go assets templates renderer viewmodel presentation
	@echo "Built $(UIKIT_BUNDLE)"
