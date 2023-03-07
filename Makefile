.PHONY: init proto test

# go env
GOPROXY     := "https://goproxy.cn,direct"
GOOS        := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
GOARCH      := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
CGO_LDFLAGS := "-static"
CC          := musl-gcc

GOENV := GO111MODULE=on
GOENV += GOPROXY=$(GOPROXY)
GOENV += CC=$(CC)
GOENV += CGO_ENABLED=1 CGO_LDFLAGS=$(CGO_LDFLAGS)
GOENV += GOOS=$(GOOS) GOARCH=$(GOARCH)

# go
GO := go

# go test
GO_TEST ?= $(GO) test

# test flags
TEST_FLAGS := -v

install_grpc_protobuf:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

proto:
	@bash mk-proto.sh

init:
	$(GO) mod init

test:
	$(GO_TEST) $(TEST_FLAGS) ./rpc/...

clean:
	rm -rf proto
