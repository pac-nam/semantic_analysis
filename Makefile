# MAIN SETTINGS

BINARY=semantic
VERSION=0.0.0
USERNAME=tbleuse

#GO SETTINGS
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOARCH=$(shell $(GOCMD) env GOARCH)

#SOFTWARE SETTINGS
COMPILE_DATE=$(shell date +"%d/%m/%Y - %H:%M")
LDFLAGS= -ldflags '-X "kdatapack/bin.VERSION_NUMBER=$(VERSION)" -X "kdatapack/bin.COMPILE_DATE=$(COMPILE_DATE)" -X "kdatapack/bin.USERNAME=$(USERNAME)"'
BINARY_NAME=$(BINARY)-$(VERSION)

#DEBIAN SETTINGS
DEB_OUTPUT=$(BINARY_NAME)/DEBIAN
DEB_BIN=$(BINARY_NAME)/usr/bin/capitaldata

all: deps build

build:
	@echo "> Building..."
	$(GOBUILD) -o $(BINARY_NAME)_standalone $(LDFLAGS) $(BUILD_TAGS)
	@echo ""

clean:
	@echo "> Cleaning folders..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)_standalone
	rm -fr $(BINARY_NAME)
	rm -f $(BINARY_NAME).deb
	@echo ""

fclean: clean
	@echo "> Cleaning potential binaries..."
	rm -rf $(BINARY)-[0-9].*
	@echo ""

deps:
	@echo "> Go getting project dependencies..."
	$(GOGET) $(BUILD_TAGS)
	@echo ""

debian: makedir control build
	@echo "> Packaging for upload to debian repository..."
	cp $(BINARY_NAME)_standalone $(DEB_BIN)/$(BINARY_NAME)
	dpkg-deb --build $(BINARY_NAME)
	@echo "======="
	@echo "> Uploading to debian repository..."
	scp $(BINARY_NAME).deb bucket@apt.kdata.fr:./incoming/
	@echo "All done."

makedir:
	mkdir -p $(DEB_BIN)
	mkdir -p $(DEB_OUTPUT)

control:
	@echo "> Generating DEBIAN/control..."
	@echo "Package: $(BINARY_NAME)" > $(DEB_OUTPUT)/control
	@echo "Version: $(VERSION)" >> $(DEB_OUTPUT)/control
	@echo "Section: base" >> $(DEB_OUTPUT)/control
	@echo "Priority: optional" >> $(DEB_OUTPUT)/control
	@echo "Architecture: $(GOARCH)" >> $(DEB_OUTPUT)/control
	@echo "Maintainer: CAPITALDATA" >> $(DEB_OUTPUT)/control
	@echo "Description: Generic Description" >> $(DEB_OUTPUT)/control
	@echo ""

re: clean all

.PHONY: all build clean deps debian makedir control re
