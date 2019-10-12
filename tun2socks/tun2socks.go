package tun2socks

import (
	"log"
	"net"
	"syscall"

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

// SetNonblock puts the fd in blocking or non-blocking mode.
func SetNonblock(fd int, nonblocking bool) bool {
	err := syscall.SetNonblock(fd, nonblocking)
	if err != nil {
		return false
	}
	return true
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
func Start(packetFlow PacketFlow, socks5Server string, fakeIPStart string, fakeIPStop string) {
	if packetFlow != nil {
		if lwipStack == nil {
			// Setup the lwIP stack.
			lwipStack = core.NewLWIPStack()
		}

		fakeDNS := fakedns.NewSimpleFakeDns(fakeIPStart, fakeIPStop)
		//var fakeDNS dns.FakeDns
		// Register tun2socks connection handlers.
		proxyAddr, err := net.ResolveTCPAddr("tcp", socks5Server)
		proxyHost := proxyAddr.IP.String()
		proxyPort := uint16(proxyAddr.Port)
		if err != nil {
			log.Fatalf("invalid proxy server address: %v", err)
		}
		core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, proxyPort, fakeDNS, nil))
		core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, proxyPort, 30, nil, fakeDNS, nil))

		// Write IP packets back to TUN.
		core.RegisterOutputFn(func(data []byte) (int, error) {
			packetFlow.WritePacket(data)
			return len(data), nil
		})

		isStopped = false
	}
}
