NAME=eks-operator
AUTHOR ?= appvia
AUTHOR_EMAIL=info@appvia.io
BUILD_TIME=$(shell date '+%s')
DEPS=$(shell go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
GIT_SHA=$(shell git --no-pager describe --always --dirty)
GOVERSION ?= 1.12.7
HARDWARE=$(shell uname -m)
LFLAGS ?= -X main.gitsha=${GIT_SHA} -X main.compiled=${BUILD_TIME}
PACKAGES=$(shell go list ./...)
REGISTRY=quay.io
ROOT_DIR=${PWD}
VERSION ?= $(shell awk '/Version.*=/ { print $$3 }' version/version.go | sed 's/"//g')
VETARGS ?= -asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -unsafeptr

.PHONY: test authors changelog build docker static release lint cover vet glide-install

default: build

golang:
	@echo "--> Go Version"
	@go version

build: golang deps operator-gen schema-gen
	@echo "--> Compiling the project"
	@mkdir -p bin
	go build -ldflags "${LFLAGS}" -o bin/${NAME} cmd/manager/*.go

static: golang deps operator-gen schema-gen
	@echo "--> Compiling the static binary"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags "-w ${LFLAGS}" -o bin/${NAME} cmd/manager/*.go

operator-gen:
	@echo "--> Generating Code via operator-sdk"
	@operator-sdk generate k8s
	@operator-sdk generate openapi

schema-gen:
	@echo "--> Generate the CRD Schema for class registration"
	@echo "--> pkg/apis/schema/schema.go"
	@go run $(GOPATH)/src/github.com/appvia/hub-apis/cmd/schema-gen/main.go \
		-crd-path ./deploy/crds \
		-package schema \
		-schema-file pkg/apis/schema/schema.go

docker-build:
	@echo "--> Compiling the project"
	docker run --rm \
		-v ${ROOT_DIR}:/go/src/github.com/${AUTHOR}/${NAME} \
		-w /go/src/github.com/${AUTHOR}/${NAME} \
		-e GOOS=linux golang:${GOVERSION} \
		make static

docker-release:
	@echo "--> Building a release image"
	@$(MAKE) static
	@$(MAKE) docker
	@docker push ${REGISTRY}/${AUTHOR}/${NAME}:${VERSION}

docker: static
	@echo "--> Building the docker image"
	operator-sdk build ${REGISTRY}/${AUTHOR}/${NAME}:${VERSION}

release: static
	mkdir -p release
	gzip -c bin/${NAME} > release/${NAME}_${VERSION}_linux_${HARDWARE}.gz
	rm -f release/${NAME}

clean:
	rm -rf ./bin 2>/dev/null
	rm -rf ./release 2>/dev/null

authors:
	@echo "--> Updating the AUTHORS"
	git log --format='%aN <%aE>' | sort -u > AUTHORS

dep-install:
	@echo "--> Installing dependencies"

deps:
	@echo "--> Installing build dependencies"

vet:
	@echo "--> Running go vet $(VETARGS) $(PACKAGES)"
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@go vet $(VETARGS) $(PACKAGES)

gofmt:
	@echo "--> Running gofmt check"
	@gofmt -s -l *.go \
	    | grep -q \.go ; if [ $$? -eq 0 ]; then \
            echo "You need to runn the make format, we have file unformatted"; \
            gofmt -s -l *.go; \
            exit 1; \
	    fi

verify:
	@echo "--> Verifying the code"
	gometalinter --disable=errcheck --disable=gocyclo --disable=gas --disable=aligncheck --errors

format:
	@echo "--> Running go fmt"
	@gofmt -s -w *.go

bench:
	@echo "--> Running go bench"
	@go test -bench=. -benchmem

coverage:
	@echo "--> Running go coverage"
	@go test -coverprofile cover.out
	@go tool cover -html=cover.out -o cover.html

lint:
	@echo "--> Running golint"
	@which golint 2>/dev/null ; if [ $$? -eq 1 ]; then \
    go get -u golang.org/x/lint/golint; \
  fi
	@golint $(PACKAGES)

cover:
	@echo "--> Running go cover"
	@go test --cover $(PACKAGES)

spelling:
	@echo "--> Checking the spelling"
	@which misspell 2>/dev/null ; if [ $$? -eq 1 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	@misspell -error *.go
	@misspell -error *.md

test:
	@echo "--> Running the tests"
	@if [ ! -d "vendor" ]; then \
		make deps; \
  fi
	@go test -v $(PACKAGES)
	@$(MAKE) golang
	@$(MAKE) gofmt
	@$(MAKE) lint
	@$(MAKE) spelling
	@$(MAKE) vet
	@$(MAKE) cover

all: test
	echo "--> Performing all tests"
	@${MAKE} verify
	@$(MAKE) bench
	@$(MAKE) coverage

changelog: release
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B > changelog
