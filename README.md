# igniter-go-libs

## How to build

### Prerequirements

* go 1.13
* gomobile https://github.com/golang/mobile
* dep https://github.com/golang/dep

### Build

```shell
GO111MODULE=off go get github.com/trojan-gfw/igniter-go-libs
cd $GOPATH/src/github.com/trojan-gfw/igniter-go-libs
dep ensure -update
make android
```
