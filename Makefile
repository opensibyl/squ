# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build_default:
	${GOCMD} build -o squ ./cmd/squ

build_macos:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 ${GOCMD} build -o squ_macos ./cmd/squ

build_linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 ${GOCMD} build -o squ_linux ./cmd/squ

build_windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 ${GOCMD} build -o squ_windows.exe ./cmd/squ

release: build_macos build_linux build_windows
