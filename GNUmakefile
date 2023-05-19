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
	golangci-lint run --timeout 2m0s ./...

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@go install $$(go list -f '{{range .Imports}}{{.}} {{end}}' tools.go)

test:
	go test -v -cpu 4 -covermode=count -coverpkg github.com/indykite/indykite-sdk-go/... -coverprofile=coverage.out.tmp ./...
	cat coverage.out.tmp | grep -v "pb.go\|pb.validate.go\|generated.go\|indykite-sdk-go/test/\|main.go\|indykite-sdk-go/examples/" > coverage.out
	rm coverage.out.tmp

cover: test
	@echo "==> generate test coverage..."
	go tool cover -html=coverage.out

upgrade:
	@echo "==> Upgrading Go"
	@GO111MODULE=on go get -u all && go mod tidy
	@echo "==> Upgrading pre-commit"
	@pre-commit autoupdate
	@echo "Please, upgrade workflows manually"

generate-proto:
	@buf generate buf.build/indykite/indykiteapis
	@go generate
	@make fmt gci
