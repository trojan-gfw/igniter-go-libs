module github.com/trojan-gfw/igniter-go-libs

go 1.16

replace github.com/eycorsican/go-tun2socks => github.com/trojan-gfw/go-tun2socks v1.16.3-0.20210309143813-f8b179dbc857

replace github.com/Dreamacro/clash => github.com/trojan-gfw/clash v0.19.1-0.20210418132805-e6ff5207426d

require (
	github.com/Dreamacro/clash v0.0.0
	github.com/djherbis/buffer v1.1.0 // indirect
	github.com/djherbis/nio v2.0.3+incompatible // indirect
	github.com/eycorsican/go-tun2socks v0.0.0
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/songgao/water v0.0.0-20200317203138-2b4b6d7c09d8
	github.com/v2pro/plz v0.0.0-20200805122259-422184e41b6e // indirect
)
