ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

UNFORMATTED=$(shell gofmt -s -l .)
#Lint
.PHONY: lint-prepare
lint-prepare:
	@echo "Preparing Linter"
	# Pin the installer to a release tag (not master) and pin the version, for
	# reproducible, checksum-verified installs. The master install.sh currently
	# fails checksum verification for every version.
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/v2.12.2/install.sh | sh -s -- -b ./bin v2.12.2

.PHONY: lint
lint:
	./bin/golangci-lint run -v ./...
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
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s v2.15.0
	./bin/gosec ./...
