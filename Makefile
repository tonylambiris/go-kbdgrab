export GO15VENDOREXPERIMENT=1

#SHA=$(shell git rev-parse --short HEAD)
#COUNT=$(shell git rev-list --count HEAD)

#BUILDTAG=${COUNT}.${SHA}
BUILDTYPE=release

all: deps bundle build

build:
	@go build -ldflags "-s -w -X main.Type=${BUILDTYPE}" -o kbdgrab

# To vendor an external dependency, run: govendor fetch path/to/repo
deps: godeps
	@echo "Fetching missing dependencies..."
	@govendor fetch +outside
	@echo "Removing unused dependencies..."
	@govendor remove +unused
	@echo "Running govendor sync..."
	@govendor sync -v

godeps:
	@echo "Installing/updating go dependencies..."
	@go get -u github.com/kardianos/govendor
	@go get -u github.com/kevinburke/go-bindata/...

bundle:
	@go-bindata LCD_Solid.ttf

clean:
	@rm -f kbdgrab bindata.go

.PHONY: build deps bundle
