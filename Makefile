export CGO_ENABLED := 1
include Makefile.Inc

test test-all: get-gpu-setup
	$(ULIMIT) CGO_LDFLAGS="$(CGO_TEST_LDFLAGS)" go test -v -timeout 0 -p 1 ./...
.PHONY: test test-all

compile-test: get-gpu-setup
	CGO_LDFLAGS="$(CGO_TEST_LDFLAGS)" go test -v -c -o $(BIN_DIR)test$(EXE) ./...
.PHONY: compile-test

ifeq ($(HOST_OS),$(filter $(HOST_OS),linux darwin))
compile-windows-test:
	CC=x86_64-w64-mingw32-gcc $(MAKE) GOOS=windows GOARCH=amd64 BIN_DIR=$(PROJ_DIR)build/ compile-test
endif
.PHONY: compile-windows-test

install: get-gpu-setup
	go mod download
	GO111MODULE=off go get golang.org/x/lint/golint
.PHONY: install

build: $(BIN_DIR)post$(EXE) #$(BIN_DIR)spacemesh-init$(EXE)
.PHONY: build

$(BIN_DIR)post$(EXE): get-gpu-setup
	go build -o $@ .

$(BIN_DIR)spacemesh-init$(EXE): get-gpu-setup
	cd cmd/init && go build -o $@ .

SHA = $(shell git rev-parse --short HEAD)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
DOCKER_IMAGE := post:$(BRANCH)

ifeq ($(BRANCH),$(filter $(BRANCH),staging trying))
  DOCKER_IMAGE = $(DOCKER_IMAGE_REPO):$(SHA)
endif

dockerbuild-go:
	docker build -t $(DOCKER_IMAGE) .
.PHONY: dockerbuild-go

test-tidy:
	# Working directory must be clean, or this test would be destructive
	git diff --quiet || (echo "\033[0;31mWorking directory not clean!\033[0m" && git --no-pager diff && exit 1)
	# We expect `go mod tidy` not to change anything, the test should fail otherwise
	make tidy
	git diff --exit-code || (git --no-pager diff && git checkout . && exit 1)
.PHONY: test-tidy

test-fmt:
	git diff --quiet || (echo "\033[0;31mWorking directory not clean!\033[0m" && git --no-pager diff && exit 1)
	# We expect `go fmt` not to change anything, the test should fail otherwise
	go fmt ./...
	git diff --exit-code || (git --no-pager diff && git checkout . && exit 1)
.PHONY: test-fmt

lint:
	golint --set_exit_status ./...
	go vet ./...
.PHONY: lint
