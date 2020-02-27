# igniter-go-libs

## How to build

### Prerequirements

* go >= 1.13
* gomobile https://github.com/golang/mobile

### Install Go

Please make sure you have at least Go 1.13 installed correctly

```shell
go version
```

You should see something like
```
go version go1.13.8 linux/amd64
```

### Install and upgrade gomobile

In this example, We will install `gomobile` and its dependencies to `$GOPATH`.

You can find the executable in `$GOPATH/bin`, to make it easy to use, we add it to `PATH`.

To upgrade, re-run the install commands.

```shell
export GOPATH=$HOME/go
go get -u -d -v golang.org/x/mobile/cmd/gomobile
go build -a -v golang.org/x/mobile/cmd/gomobile
go install -v golang.org/x/mobile/cmd/gomobile
# Make sure we can execute gomobile and gobind directly
export PATH="$GOPATH/bin:$PATH"
# prepare gomobile component: gobind
gomobile init
```

### Build

```shell
# If you cannot execute command, make sure your PATH is correct

# Point to your Android SDK root
# change to your installation configuration please
export ANDROID_HOME=/path/to/your/android/sdk

# Clone this repository
pushd /path/to/git/repository/just/cloned

make android
```

### Development Guide

#### How to switch between local and remove module

If you want to say, point to the local version of a dependency in Go rather than the one over the web, use the replace keyword.

The replace line goes above your require statements, like so:

```
module github.com/person/foo

replace github.com/person/bar => /Users/person/Projects/bar

require (
	github.com/person/bar v1.0.0
)

```

And now when you compile this module (go install), it will use your local code rather than the other dependency.

According to the docs, you do need to make sure that the code you’re pointing to also has a go.mod file.

```shell
go mod edit -replace github.com/eycorsican/go-tun2socks=/path/to/my/local/github/igniter-deps/go-tun2socks
```

It looks everything is fine, we can push changes in your local repository to the remote one.

In this example, the commit we just pushed is `efeee82`.

To switch to the version over the web, a.k.a. the commit `efeee82` from module `github.com/trojan-gfw/go-tun2socks`

```shell
go mod edit -replace github.com/eycorsican/go-tun2socks=github.com/trojan-gfw/go-tun2socks@efeee82
```

To check which dependency you are truely using:

```shell
go list -m all
```

#### How to make C/C++ source changes take effect

For go tools, e.g. `gobind`, it only check `*.go` file changes.

Please make sure to change at least one golang source by adding whitespaces or line breaks or whatever you like.

The go tools notify `*.go` changes, it will recompile the C/C++ sources.

