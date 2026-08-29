APP_NAME := govnoed-user-items
BIN_DIR := bin
CMD_PATH := ./cmd/api

.PHONY: all \
	build \
	windows-amd64 \
	windows-arm64 \
	windows-386 \
	linux-amd64 \
	linux-arm64 \
	linux-386 \
	check \
	run \
	clean

all: build

build:
	$(MAKE) windows-amd64
	$(MAKE) linux-amd64

windows-amd64:
	GOOS=windows GOARCH=amd64 go build \
		-o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe \
		$(CMD_PATH)

windows-arm64:
	GOOS=windows GOARCH=arm64 go build \
		-o $(BIN_DIR)/$(APP_NAME)-windows-arm64.exe \
		$(CMD_PATH)

windows-386:
	GOOS=windows GOARCH=386 go build \
		-o $(BIN_DIR)/$(APP_NAME)-windows-386.exe \
		$(CMD_PATH)

linux-amd64:
	GOOS=linux GOARCH=amd64 go build \
		-o $(BIN_DIR)/$(APP_NAME)-linux-amd64 \
		$(CMD_PATH)

linux-arm64:
	GOOS=linux GOARCH=arm64 go build \
		-o $(BIN_DIR)/$(APP_NAME)-linux-arm64 \
		$(CMD_PATH)

linux-386:
	GOOS=linux GOARCH=386 go build \
		-o $(BIN_DIR)/$(APP_NAME)-linux-386 \
		$(CMD_PATH)

check:
	go vet ./...

run:
	go run $(CMD_PATH)

clean:
	rm -f $(BIN_DIR)/$(APP_NAME)-*