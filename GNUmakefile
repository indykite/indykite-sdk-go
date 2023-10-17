.PHONY: test

default:

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w .

goimports: gci

gci:
	@echo "==> Fixing imports code with gci..."
	gci write -s standard -s default -s "prefix(github.com/indykite/indykite-sdk-go)" -s blank -s dot .

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run --timeout 3m0s ./...

lint_fix:
	@echo "==> Fixing source code against linters..."
	@golangci-lint run --fix --timeout 3m0s ./...

fieldalignment:
	@echo "==> Analyzer structs and rearranged to use less memory with fieldalignment..."
	fieldalignment -fix -test=false ./...

install-tools:
	@echo Installing tools
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install github.com/daixiang0/gci@latest
	@go install go.uber.org/mock/mockgen@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
	@go install github.com/onsi/ginkgo/v2/ginkgo@latest
	@go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	@echo Installation completed

test:
	go test -v -cpu 4 -covermode=count -coverpkg github.com/indykite/indykite-sdk-go/... -coverprofile=coverage.out.tmp ./...
	cat coverage.out.tmp | grep -v "pb.go\|pb.validate.go\|generated.go\|indykite-sdk-go/test/\|main.go\|indykite-sdk-go/examples/" > coverage.out
	rm coverage.out.tmp

cover: test
	@echo "==> generate test coverage..."
	go tool cover -html=coverage.out

upgrade: upgrade_deps generate-proto

upgrade_deps:
	@echo "==> Upgrading Go"
	@GO111MODULE=on go get -u all && go mod tidy
	@echo "==> Upgrading pre-commit"
	@pre-commit autoupdate
	@echo "Please, upgrade workflows manually"

generate-proto:
	@echo "==> Generate Proto messages"
	@buf generate buf.build/indykite/indykiteapis
	@go generate
	@make fmt gci
