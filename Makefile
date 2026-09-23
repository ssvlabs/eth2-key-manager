ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

# Tools (golangci-lint, ...) are pinned in tool.mod and run via the Go toolchain,
# so their versions are reproducible and checksum-verified through tool.sum -- no
# curl|sh installer needed. Requires a Go toolchain >= the version in tool.mod.
RUN_TOOL=go tool -modfile=tool.mod

#Lint (incl. gosec and gofmt, see .golangci.yml)
.PHONY: lint
lint:
	GOWORK=off $(RUN_TOOL) github.com/golangci/golangci-lint/v2/cmd/golangci-lint run -v ./...

#Test
.PHONY: full-test
full-test:
	@echo "Running the full test..."
	@go test -tags blst_enabled -timeout 20m ${COV_CMD} -race -p 1 -v ./...
