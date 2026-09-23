# Tools are pinned in tool.mod (checksum-verified via tool.sum) and run with
# `go tool`, which needs a Go toolchain >= tool.mod's go directive. See tool.mod
# for how to bump them. GOWORK=off because -modfile fails in workspace mode
# (e.g. with a go.work anywhere up the tree).
RUN_TOOL=GOWORK=off go tool -modfile=tool.mod

#Lint (incl. gosec and gofmt, see .golangci.yml)
.PHONY: lint
lint:
	$(RUN_TOOL) github.com/golangci/golangci-lint/v2/cmd/golangci-lint run -v ./...

#Test
.PHONY: full-test
full-test:
	@echo "Running the full test..."
	@go test -tags blst_enabled -timeout 20m ${COV_CMD} -race -p 1 -count=1 -v ./...
