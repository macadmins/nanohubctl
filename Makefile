VERSION = $(shell git describe --tags --always)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"
OSARCH=$(shell go env GOHOSTOS)-$(shell go env GOHOSTARCH)

DDMCTL=\
	ddmctl-darwin-amd64 \
	ddmctl-darwin-arm64 \
	ddmctl-linux-amd64 \
	ddmctl-linux-arm64 \
	ddmctl-linux-arm \
	ddmctl-windows-amd64.exe

my: ddmctl-$(OSARCH)

$(DDMCTL): main.go
	GOOS=$(word 2,$(subst -, ,$@)) GOARCH=$(word 3,$(subst -, ,$(subst .exe,,$@))) go build $(LDFLAGS) -o $@ ./$<

ddmctl-%-$(VERSION).zip: ddmctl-%.exe
	rm -rf $(subst .zip,,$@)
	mkdir $(subst .zip,,$@)
	ln $^ $(subst .zip,,$@)
	zip -r $@ $(subst .zip,,$@)
	rm -rf $(subst .zip,,$@)

ddmctl-%-$(VERSION).zip: ddmctl-%
	rm -rf $(subst .zip,,$@)
	mkdir $(subst .zip,,$@)
	ln $^ $(subst .zip,,$@)
	zip -r $@ $(subst .zip,,$@)
	rm -rf $(subst .zip,,$@)

clean:
	rm -rf ddmctl-*

release: $(foreach bin,$(DDMCTL),$(subst .exe,,$(bin))-$(VERSION).zip)

test:
	go test -v -cover -race ./...

.PHONY: my $(DDMCTL) clean release test
