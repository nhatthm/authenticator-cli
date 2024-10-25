APP ?= authenticator

WORK_DIR = $(shell pwd)
BUILD_DIR ?= out
VENDOR_DIR = vendor

GOLANGCI_LINT_VERSION ?= v1.61.0

GO ?= go
GOLANGCI_LINT ?= $(shell go env GOPATH)/bin/golangci-lint-$(GOLANGCI_LINT_VERSION)

VHS_TAPES = $(subst .tape,.gif,$(shell find resources/docs -type f -name "*.tape" | sort))

.PHONY: $(VENDOR_DIR)
$(VENDOR_DIR):
	@mkdir -p $(VENDOR_DIR)
	@$(GO) mod vendor
	@$(GO) mod tidy

.PHONY: generate
generate: generate-docs

.PHONY: generate-docs
generate-docs: $(VHS_TAPES)

.PHONY: $(VHS_TAPES)
$(VHS_TAPES):
	@echo ">> generate $@"
	@cd resources/docs; \
		env PATH="$(WORK_DIR)/out:$$PATH" vhs < $(notdir $(basename $@)).tape

.PHONY: lint
lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run

.PHONY: test
test: test-unit

## Run unit tests
.PHONY: test-unit
test-unit:
	@echo ">> unit test"
	@$(GO) test -gcflags=-l -coverprofile=unit.coverprofile -covermode=atomic -race ./...

#.PHONY: test-integration
#test-integration:
#	@echo ">> integration test"
#	@$(GO) test ./features/... -gcflags=-l -coverprofile=features.coverprofile -coverpkg ./... -race --godog

.PHONY: build
build:
	@$(GO) build -ldflags "$(shell ./resources/scripts/build_args)" -o $(BUILD_DIR)/$(APP) cmd/$(APP)/main.go && \
		chmod +x $(BUILD_DIR)/$(APP)

.PHONY: $(GITHUB_ENV)
$(GITHUB_ENV):
	@echo "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)" >>"$@"

.PHONY: $(GITHUB_OUTPUT)
$(GITHUB_OUTPUT):
	@echo "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)" >>"$@"

$(GOLANGCI_LINT):
	@echo "$(OK_COLOR)==> Installing golangci-lint $(GOLANGCI_LINT_VERSION)$(NO_COLOR)"; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin "$(GOLANGCI_LINT_VERSION)"
	@mv ./bin/golangci-lint $(GOLANGCI_LINT)
