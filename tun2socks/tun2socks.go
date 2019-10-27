package tun2socks

import (
	"log"
	"net"

	"github.com/eycorsican/go-tun2socks/common/dns/cache"
	"github.com/eycorsican/go-tun2socks/common/dns/fakedns"
	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
)

var lwipStack core.LWIPStack
var isStopped = false

// PacketFlow should be implemented in Java/Kotlin.
type PacketFlow interface {
	// WritePacket should writes packets to the TUN fd.
	WritePacket(packet []byte)
}

// InputPacket Write IP packets to the lwIP stack. Call this function in the main loop of
// the VpnService in Java/Kotlin, which should reads packets from the TUN fd.
func InputPacket(data []byte) {
	if lwipStack != nil {
		lwipStack.Write(data)
	}
}

// Stop stop it
func Stop() {
	isStopped = true
	if lwipStack != nil {
		lwipStack.Close()
		lwipStack = nil
	}
}

// Start sets up lwIP stack, starts a Tun2socks instance
func Start(packetFlow PacketFlow, socks5Server string, fakeIPStart string, fakeIPStop string) int {
	if packetFlow != nil {
		if lwipStack == nil {
			// Setup the lwIP stack.
			lwipStack = core.NewLWIPStack()
		}

		// Register tun2socks connection handlers.
		proxyAddr, err := net.ResolveTCPAddr("tcp", socks5Server)
		proxyHost := proxyAddr.IP.String()
		proxyPort := uint16(proxyAddr.Port)
		if err != nil {
			log.Fatalf("invalid proxy server address: %v", err)
			return -1
		}
		cacheDNS := cache.NewSimpleDnsCache()
		if fakeIPStart != "" && fakeIPStop != "" {
			fakeDNS := fakedns.NewSimpleFakeDns(fakeIPStart, fakeIPStop)
			core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, proxyPort, fakeDNS, nil))
			core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, proxyPort, 30, cacheDNS, fakeDNS, nil))
		} else {
			core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, proxyPort, nil, nil))
			core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, proxyPort, 30, cacheDNS, nil, nil))
		}

		// Write IP packets back to TUN.
		core.RegisterOutputFn(func(data []byte) (int, error) {
			packetFlow.WritePacket(data)
			return len(data), nil
		})

		isStopped = false
	}
	return 0
}
