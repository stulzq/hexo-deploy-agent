LDFLAGS := -s -w
VERSION := "v0.1.0"
APPNAME := "hexo-deploy-agent"

build: clean build-linux-amd64

build-linux-amd64:
	@rm -rf bin/output
	@GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/output/hexo-deploy-agent ./cmd
	@mkdir bin/output/conf/
	@cp conf/config.yml bin/output/conf/
	@cd bin && tar -zvcf $(APPNAME)-$(VERSION)-linux-amd64.tar.gz output
	@echo "linux-amd64 build and package"

run:
	@go run ./cmd

clean:
	@rm -rf bin/*

.PHONY: run build