########################### env ##############################
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
export GOSUMDB=sum.golang.org
##############################################################
-include protobuf/Makefile
GOPATH:=$(shell go env GOPATH)
SRC=$(shell find . -name "*.go")
VERSION=$(shell git describe --always --abbrev=10)
API_PROTO_FILES=$(shell find api -name "*.proto")
INJECT_PB_GO=$(shell find api -name "*.pb.go" | grep -v "_")
# golangci-lint
LINTER := golangci-lint

$(LINTER):
	@curl -SL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

.PHONY: init
# init
init:
	$(info ******************** init ********************)
	@test -e configs/config.yaml || cp configs/config.yaml.bak configs/config.yaml
	go install github.com/google/wire/cmd/wire@latest

.PHONY: fmt
# go fmt
fmt:
	$(info ******************** checking formatting ********************)
	@go fmt ./...

.PHONY: lint
# lint
lint: fmt
	$(info ******************** running lint ********************)
	$(LINTER) run ./... --timeout=10m

.PHONY: test
# go test
test:
	$(info ******************** running testing ********************)
	@go test -v -cover -gcflags=all=-l -coverprofile=coverage.out ./...

.PHONY: build
# go build
build:
	$(info ******************** go build ********************)
	@mkdir -p bin/ && CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -mod=vendor -o ./bin/ ./...


.PHONY: wire
# generate wire
wire:
	$(info ******************** wire ********************)
	find cmd -type d -depth 1 -print | xargs -L 1 bash -c 'cd "$$0" && pwd && wire'

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
    helpMessage = match(lastLine, /^# (.*)/); \
        if (helpMessage) { \
            helpCommand = substr($$1, 0, index($$1, ":")-1); \
            helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
            printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
        } \
    } \
    { lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
