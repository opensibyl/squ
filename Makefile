# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build_macos:
	GOOS=darwin GOARCH=amd64 ${GOCMD} build -o squ_macos ./cmd/squ

build_linux:
	GOOS=linux GOARCH=amd64 ${GOCMD} build -o squ_linux ./cmd/squ

build_windows:
	GOOS=windows GOARCH=amd64 ${GOCMD} build -o squ_windows.exe ./cmd/squ
