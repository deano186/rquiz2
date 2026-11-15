.DEFAULT_GOAL := all

.PHONY: all clean fmt vet build android windows

# Build every version: Linux, Android and Windows
all: build android windows

clean:
	go clean ./...
	rm -f rquiz2-android rquiz2.exe
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build

# Static arm64 binary for Android (e.g. run under Termux)
android: vet
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o rquiz2-android .

# 64-bit Windows console executable
windows: vet
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o rquiz2.exe .
