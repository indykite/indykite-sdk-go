.PHONY: default build vet fmt goimports gci lint lint_fix fieldalignment test test_race integration report fixtures fixtures_destroy cover tidy upgrade install-tools

default: build

build:
	@echo "==> Building..."
	go build ./...

vet:
	@echo "==> go vet..."
	go vet ./...

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w .

goimports: gci

gci:
	@echo "==> Fixing imports code with gci..."
	gci write -s standard -s default -s "prefix(github.com/indykite/indykite-sdk-go)" -s blank -s dot .

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run --path-mode=abs --timeout 3m0s ./...

lint_fix:
	@echo "==> Fixing source code against linters..."
	@golangci-lint run --fix --timeout 3m0s ./...

fieldalignment:
	@echo "==> Rearranging struct fields to use less memory with fieldalignment..."
	fieldalignment -fix -test=false ./...

test:
	go test -cpu 4 -covermode=count -coverpkg github.com/indykite/indykite-sdk-go/... -coverprofile=coverage.out.tmp ./...
	cat coverage.out.tmp | grep -v "indykite-sdk-go/agents/\|indykite-sdk-go/test/" > coverage.out
	rm coverage.out.tmp

test_race:
	go test -race -count=1 ./...

# -p 1 serializes the per-package test binaries: they all hit the same live
# environment, and concurrent graph writes make fixture queries flaky.
integration:
	go test -tags integration -count=1 -p 1 -run TestIntegration ./...

# Provision / remove the integration-test fixtures (see test/README.md).
fixtures:
	go run ./test/setup apply

fixtures_destroy:
	go run ./test/setup destroy

# report runs the integration suite once and renders test-report.html; the
# go-sdk-tests image (docker/infra) uploads it to the results bucket. bash +
# pipefail so a test failure is the target's exit code, not go-test-report's.
report:
	bash -o pipefail -c 'go test -v -json -tags integration -count=1 -p 1 -run TestIntegration ./... \
		| go-test-report -o test-report.html -t "Go SDK Tests report"'

cover: test
	@echo "==> generate test coverage..."
	go tool cover -html=coverage.out

tidy:
	go mod tidy

upgrade:
	@echo "==> Upgrading Go dependencies"
	@GO111MODULE=on go get -u all && go mod tidy
	@echo "==> Upgrading pre-commit"
	@pre-commit autoupdate

install-tools:
	@echo Installing tools
	@go install github.com/daixiang0/gci@latest
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	@go install github.com/vakenbolt/go-test-report@latest
	@echo Installation completed
