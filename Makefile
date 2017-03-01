export GO15VENDOREXPERIMENT=1

#SHA=$(shell git rev-parse --short HEAD)
#COUNT=$(shell git rev-list --count HEAD)

#BUILDTAG=${COUNT}.${SHA}
BUILDTYPE=release

all: deps bundle build

build:
	@go build -ldflags "-s -w -X main.Type=${BUILDTYPE}" -o kbdgrab

deps:
	@go get github.com/BurntSushi/xgbutil/...
	@go get github.com/jteeuwen/go-bindata/...

bundle:
	@go-bindata LCD_Solid.ttf

clean:
	@rm -f kbdgrab bindata.go

.PHONY: build deps bundle
