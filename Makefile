ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

UNFORMATTED=$(shell gofmt -s -l .)

# Tools (golangci-lint, ...) are pinned in tool.mod and run via the Go toolchain,
# so their versions are reproducible and checksum-verified through tool.sum -- no
# curl|sh installer needed. Requires a Go toolchain >= the version in tool.mod.
RUN_TOOL=go tool -modfile=tool.mod

#Lint
.PHONY: lint
lint:
	GOWORK=off $(RUN_TOOL) github.com/golangci/golangci-lint/v2/cmd/golangci-lint run -v ./...
	@echo "Checking for unformatted files"
	if [ ! -z "${UNFORMATTED}" ]; then \
		echo "The following files are not formatted: \n${UNFORMATTED}"; \
	fi

#Test
.PHONY: full-test
full-test:
	@echo "Running the full test..."
	@go test -tags blst_enabled -timeout 20m ${COV_CMD} -race -p 1 -v ./...

#Gosec
.PHONY: gosec
gosec:
	@echo "Running the gosec check"
	# Pin the installer to a release tag (not master) for a reproducible install.
	curl -sfL https://raw.githubusercontent.com/securego/gosec/v2.15.0/install.sh | sh -s -- -b ./bin v2.15.0
	./bin/gosec ./...
