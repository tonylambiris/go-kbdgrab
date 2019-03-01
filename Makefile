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

build:
	@go build -ldflags "-s -w -X main.Build=${BUILDTAG} -X main.Type=${BUILDTYPE}" \
		-gcflags="all=-trimpath=${GOPATH}/src" \
		-asmflags="all=-trimpath=${GOPATH}/src" \
		-o kbdgrab

# To vendor an external dependency, run: dep -add path/to/repo
deps: godeps
	@echo "Running dependency check..."
	@dep check

godeps:
	@echo "Installing/updating go dependencies..."
	@go get -v github.com/jteeuwen/go-bindata/...

update:
	@echo "Updating vendored dependencies..."
	@dep ensure -update

bundle:
	@go-bindata LCD_Solid.ttf

clean:
	@rm -f kbdgrab bindata.go

.PHONY: build deps bundle
