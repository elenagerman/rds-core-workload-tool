# Export GO111MODULE=on to enable project to be built from within GOPATH/src
export GO111MODULE=on

.PHONY: client-bin \
		client-pod

deps-update:
	go mod tidy && \
	go mod vendor

gofmt:
	@echo "Running gofmt"
	gofmt -s -l `find . -path ./vendor -prune -o -type f -name '*.go' -print`

build:
	@echo "Making testcmd binary"
	scripts/build-testcmd-bin.sh

lint:
	@echo "Running go lint"
	scripts/golangci-lint.sh
