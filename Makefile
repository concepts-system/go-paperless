VERSION		:=	$(shell cat ./VERSION)
BUILD_DATE	:=	$(shell date +%Y-%m-%d\ %H:%M)
GIT_REF		:=	$(shell git rev-parse --abbrev-ref HEAD)
BINARY_NAME	:=	go-paperless

HAS_RICHGO := $(shell which richgo)
ifdef HAS_RICHGO
    GOTEST=richgo test
else
    GOTEST=go test
endif

all: install

clean:
	rm -f $(BINARY_NAME)

install:
	go install

update-dependencies:
    go get -u -t ./...
    make clean-dependencies

clean-dependencies:
    go mod tidy

lint:
    golangci-lint run

test:
	$(GOTEST) ./...

format:
	go fmt ./...

build:
	go build -o $(BINARY_NAME) main.go

build-docker:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.release=true' -X 'main.buildDate=$(BUILD_DATE)' -X 'main.version=$(VERSION)'" -o $(BINARY_NAME) main.go

build-release:
	go build -ldflags "-X 'main.release=true' -X 'main.buildDate=$(BUILD_DATE)' -X 'main.version=$(VERSION)'" -o $(BINARY_NAME) main.go

release:
	git tag -a "v$(VERSION)" -m "Release v$(VERSION)" || true
	git push origin "v$(VERSION)"

.PHONY: clean install test fmt build image image-release release lint update-dependencies clean-dependencies
