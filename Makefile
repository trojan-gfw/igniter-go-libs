BUILDDIR=$(shell pwd)/build
IMPORT_PATH= \
	github.com/trojan-gfw/igniter-go-libs/clash \
	github.com/trojan-gfw/igniter-go-libs/tun2socks \
	github.com/trojan-gfw/igniter-go-libs/freeport \
	github.com/trojan-gfw/igniter-go-libs/util

all: ios android

ios: clean
	mkdir -p $(BUILDDIR)
	gomobile bind -o $(BUILDDIR)/golibs.framework -a -ldflags '-w' -target=ios $(IMPORT_PATH)


android: clean
	mkdir -p $(BUILDDIR)
	env GO111MODULE="on" gomobile bind -o $(BUILDDIR)/golibs.aar -a -v -x -androidapi 23 -ldflags '-w' -target=android $(IMPORT_PATH)

clean:
	gomobile clean
	rm -rf $(BUILDDIR)

cleanmodcache:
	go clean -modcache
