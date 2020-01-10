BUILDDIR=$(shell pwd)/build
SOURCES= \
	github.com/trojan-gfw/igniter-go-libs/clash \
	github.com/trojan-gfw/igniter-go-libs/tun2socks \
	github.com/trojan-gfw/igniter-go-libs/freeport

all: ios android

ios: clean
	mkdir -p $(BUILDDIR)
	gomobile bind -o $(BUILDDIR)/golibs.framework -a -ldflags '-s -w' -target=ios $(SOURCES)


android: clean
	mkdir -p $(BUILDDIR)
	gomobile bind -o $(BUILDDIR)/golibs.aar -a -ldflags '-s -w' -target=android $(SOURCES)

clean:
	rm -rf $(BUILDDIR)
