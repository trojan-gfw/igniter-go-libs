package tun2socks

import (
	"fmt"
	"strings"

	"github.com/xjasonlyu/tun2socks/v2/engine"
	"github.com/xjasonlyu/tun2socks/v2/log"
)

type Tun2socksStartOptions struct {
	TunFd        int
	Socks5Server string
	FakeIPRange  string
	MTU          int
	EnableIPv6   bool
	AllowLan     bool
}

var (
	mtuUsed        int
	cachedLoglevel string = "info"
	cachedKey             = new(engine.Key)
)

// Stop stop it
func Stop() {
	log.Infof("enter stop")
	log.Infof("this is no-op when using getFd() from Android java side, we do NOT own the tun fd")
	/*
		if err := engine.Stop(); err != nil {
			log.Errorf("Failed to stop engine %v", err)
		}
	*/
}

// Start sets up lwIP stack, starts a Tun2socks instance
func Start(opt *Tun2socksStartOptions) int {
	cachedKey.MTU = opt.MTU
	cachedKey.UDPTimeout = 600
	cachedKey.Device = fmt.Sprintf("fd://%d", opt.TunFd)
	cachedKey.LogLevel = cachedLoglevel
	cachedKey.Proxy = fmt.Sprintf("socks5://%s", opt.Socks5Server)
	cachedKey.EnableIPv6 = opt.EnableIPv6
	engine.Insert(cachedKey)

	engine.Start()
	log.Infof("Running tun2socks")

	return 0
}

// SetLoglevel set tun2socks log level
// possible input: debug/info/warn/error/none
// Log level [debug|info|warning|error|silent] for internals
func SetLoglevel(logLevel string) {
	// Set log level.
	switch strings.ToLower(logLevel) {
	case "debug":
		cachedLoglevel = "debug"
	case "info":
		cachedLoglevel = "info"
	case "warn":
		cachedLoglevel = "warning"
	case "error":
		cachedLoglevel = "error"
	case "none":
		cachedLoglevel = "silent"
	default:
		panic("unsupport logging level")
	}
	log.Infof("cached LogLevel: %v", cachedLoglevel)
}
