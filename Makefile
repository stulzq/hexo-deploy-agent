LDFLAGS := -s -w
BIN_NAME := "hexo_deploy_agent"

build:
	@env CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BIN_NAME) ./cmd

run:
	@go run ./cmd

build-linux:
	@GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/linux-amd64/$(BIN_NAME) ./cmd

build-darwin:
	@GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/darwin-amd64/$(BIN_NAME) ./cmd

clean:
	@rm -rf bin/*

.PHONY: run build build-linux build-darwin
