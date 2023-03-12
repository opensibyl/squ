# General
WORKDIR = $(PWD)

# Go parameters
GOCMD = go
GOTEST = $(GOCMD) test

build_all:
	${GOCMD} build ./cmd/squ

test:
	$(GOTEST) ./...
