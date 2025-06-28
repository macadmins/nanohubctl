VERSION = $(shell git describe --tags --always)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"
OSARCH=$(shell go env GOHOSTOS)-$(shell go env GOHOSTARCH)

NANOHUBCTL=\
	nanohubctl-darwin-amd64 \
	nanohubctl-darwin-arm64 \
	nanohubctl-linux-amd64 \
	nanohubctl-linux-arm64 \
	nanohubctl-linux-arm \
	nanohubctl-windows-amd64.exe

my: nanohubctl-$(OSARCH)

$(NANOHUBCTL): main.go
	GOOS=$(word 2,$(subst -, ,$@)) GOARCH=$(word 3,$(subst -, ,$(subst .exe,,$@))) go build $(LDFLAGS) -o $@ ./$<

nanohubctl-%-$(VERSION).zip: nanohubctl-%.exe
	rm -rf $(subst .zip,,$@)
	mkdir $(subst .zip,,$@)
	ln $^ $(subst .zip,,$@)
	zip -r $@ $(subst .zip,,$@)
	rm -rf $(subst .zip,,$@)

nanohubctl-%-$(VERSION).zip: nanohubctl-%
	rm -rf $(subst .zip,,$@)
	mkdir $(subst .zip,,$@)
	ln $^ $(subst .zip,,$@)
	zip -r $@ $(subst .zip,,$@)
	rm -rf $(subst .zip,,$@)

clean:
	rm -rf nanohubctl-*

release: $(foreach bin,$(NANOHUBCTL),$(subst .exe,,$(bin))-$(VERSION).zip)

test:
	go test -v -cover -race ./...

.PHONY: my $(NANOHUBCTL) clean release test
