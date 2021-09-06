export GO15VENDOREXPERIMENT=1

SHA=$(shell git rev-parse --short HEAD)
COUNT=$(shell git rev-list --count HEAD)

BUILDTAG=${COUNT}.${SHA}

BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
ifeq ($(BRANCH),master)
BUILDTYPE=release
else
BUILDTYPE=$(BRANCH)
endif

all: deps bundle build

build: bundle
	@go build -trimpath -ldflags \
		"-s -w -X main.Build=${BUILDTAG} -X main.Type=${BUILDTYPE}" \
		-o kbdgrab

bundle: deps
	@go-bindata LCD_Solid.ttf

clean:
	@rm -f kbdgrab bindata.go

tidy:
	@echo "Tidying up dependencies..."
	@go mod tidy

deps:
	@echo "Getting required dependencies..."
	@go get -u github.com/kevinburke/go-bindata/...

.PHONY: build deps bundle
