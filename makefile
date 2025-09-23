PHONY:
GO111MODULE=on

default:

run:
	CONFIG_PATH=config.yaml go run ./cmd/bookshop/main.go

fmt:
	@gofmt -s -w $$(go list -f "{{.Dir}}" ./...)

gci:
	@gci write -s standard -s default -s "prefix(github.com/CanobbioE/strict-clean-arch-go-webservice)" -s blank -s dot ./cmd .

lint-all:
	@golangci-lint run --timeout 2m0s ./...

lint:
	@golangci-lint run --new-from-rev=$$(git merge-base HEAD main) --timeout 2m0s ./...

test:
	go test $(shell go list ./... | grep -vE 'test|vendor') -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

generate-mocks: _generate-mocks fmt gci

_generate-mocks:
	cd internal/test && go generate

generate-proto: _proto fmt gci

_proto:
	@buf lint
	@buf format --write
	@buf generate

install-tools:
	@echo Installing tools
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install github.com/bufbuild/buf/cmd/protoc-gen-buf-lint@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/daixiang0/gci@latest
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@go install github.com/golang/mock/mockgen@v1.6.0
	@echo Installation completed



