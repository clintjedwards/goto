SEMVER = v1.0.0
GO_LDFLAGS = '-X "main.version=$(SEMVER)"'
BUILD_PATH = /tmp/test

## build: run tests and compile full app in production mode
build:
	go mod tidy
	go build -ldflags $(GO_LDFLAGS) -o $(BUILD_PATH)
