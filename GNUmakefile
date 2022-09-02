GO111MODULE=on
default:

goimports:
	@echo "==> Fixing imports code with goimports in Jarvis..."
	@goimports -local "github.com/indykite/jarvis-sdk-go" -w .

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run --timeout 2m0s ./...

download:
	@echo Download go.mod dependencies
	@go mod download

install-tools: download
	@echo Installing tools from tools.go
	@go list -f '{{range .Imports}}{{.}} {{end}}' tools.go | xargs go install

.PHONY: test
test:
	go test -v -cpu 4 -covermode=count -coverpkg github.com/indykite/jarvis-sdk-go/... -coverprofile=coverage.out.tmp ./...
	cat coverage.out.tmp | grep -v "pb.go\|pb.validate.go\|generated.go\|jarvis-sdk-go/test/\|main.go\|jarvis-sdk-go/examples/" > coverage.out
	rm coverage.out.tmp

cover: test
	@echo "==> generate test coverage..."
	go tool cover -html=coverage.out

generate-proto:
	@buf generate buf.build/indykite/indykiteapis
	@go generate
	@make goimports
