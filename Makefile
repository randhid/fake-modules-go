
BIN_OUTPUT_PATH = bin/$(shell uname -s)-$(shell uname -m)
TOOL_BIN = bin/gotools/$(shell uname -s)-$(shell uname -m)
UNAME_S ?= $(shell uname -s)a
UNAME_M ?= $(shell uname -m)

build:
	rm -f $(BIN_OUTPUT_PATH)/fake-modules-go
	go build $(LDFLAGS) -o $(BIN_OUTPUT_PATH)/fake-modules-go main.go

module.tar.gz: build
	rm -f $(BIN_OUTPUT_PATH)/module.tar.gz
	tar czf $(BIN_OUTPUT_PATH)/module.tar.gz $(BIN_OUTPUT_PATH)/fake-modules-go

fake-modules-go: *.go 
	go build -o fake-modules-go *.go

clean:
	rm -rf $(BIN_OUTPUT_PATH)/fake-modules-go $(BIN_OUTPUT_PATH)/module.tar.gz fake-modules-go

gofmt:
	gofmt -w -s .

lint: gofmt
	go mod tidy

update-rdk:
	go get go.viam.com/rdk@latest
	go mod tidy