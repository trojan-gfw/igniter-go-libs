BUILDDIR=$(shell pwd)/build
SOURCES= \
	github.com/trojan-gfw/igniter-go-libs/clash \
	github.com/trojan-gfw/igniter-go-libs/tun2socks \
	github.com/trojan-gfw/igniter-go-libs/freeport \
	github.com/trojan-gfw/igniter-go-libs/util

# pass a single dollar sign to shell
CURRENT_GOPATH="$(shell echo $$GOPATH)"

all: ios android

ios: clean
	mkdir -p $(BUILDDIR)
	gomobile bind -o $(BUILDDIR)/golibs.framework -a -ldflags '-s -w' -target=ios $(SOURCES)


android: clean
	mkdir -p $(BUILDDIR)
	env GO111MODULE="on" gomobile bind -o $(BUILDDIR)/golibs.aar -a -v -x -ldflags '-s -w' -target=android  -gcflags=-trimpath="$(CURRENT_GOPATH)" $(SOURCES)

clean:
	gomobile clean
	rm -rf $(BUILDDIR)

cleanmodcache:
	go clean -modcache
