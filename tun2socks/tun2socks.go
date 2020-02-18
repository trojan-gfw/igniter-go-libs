package tun2socks

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/eycorsican/go-tun2socks/common/dns/cache"
	"github.com/eycorsican/go-tun2socks/common/dns/fakedns"
	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
	"github.com/songgao/water"
)

var (
	lwipWriter        io.Writer
	lwipStack         core.LWIPStack
	mtuUsed           int
	stopSignalChannel chan bool
	stopReplyChannel  chan bool
	tunDev            *water.Interface
)

// Stop stop it
func Stop() {
	log.Printf("enter stop")
	log.Printf("begin close tun")
	err := tunDev.Close()
	if err != nil {
		log.Printf("close tun: %v", err)
	}
	log.Printf("send stop sig")
	close(stopSignalChannel)
	log.Printf("stop sig sent")
	<-stopReplyChannel
	if lwipStack != nil {
		log.Printf("begin close lwipstack")
		lwipStack.Close()
		lwipStack = nil
	}
}

// hack to receive tunfd
func openTunDevice(tunFd int) (*water.Interface, error) {
	file := os.NewFile(uintptr(tunFd), "tun") // dummy file path name since we already got the fd
	tunDev = &water.Interface{
		ReadWriteCloser: file,
	}
	return tunDev, nil
}

// DataPipeWorker generator
func createDataPipeWorker() chan bool {
	// a stop signal channel
	c := make(chan bool)

	// Copy packets from tun device to lwip stack, it's the main loop.
	go func(c <-chan bool) {
		var ok bool
	Loop:
		for {
			select {
			case _, ok = <-c:
				if !ok {
					log.Printf("got DataPipe stop signal")
					break Loop
				}

			default:
				// tun -> lwip
				_, err := io.CopyBuffer(lwipWriter, tunDev, make([]byte, mtuUsed))
				if err != nil {
					log.Printf("copying data failed: %v", err)
				}
			}

		}
		log.Printf("exit DataPipe loop")
		close(stopReplyChannel)
	}(c)

	return c
}

// Start sets up lwIP stack, starts a Tun2socks instance
func Start(tunFd int, socks5Server string, fakeIPStart string, fakeIPStop string, mtu int) int {

	mtuUsed = mtu
	var err error
	tunDev, err = openTunDevice(tunFd)
	if err != nil {
		log.Fatalf("failed to open tun device: %v", err)
	}

	if lwipStack == nil {
		// Setup the lwIP stack.
		lwipStack = core.NewLWIPStack()
		lwipWriter = lwipStack.(io.Writer)
	}

	// Register tun2socks connection handlers.
	proxyAddr, err := net.ResolveTCPAddr("tcp", socks5Server)
	proxyHost := proxyAddr.IP.String()
	proxyPort := uint16(proxyAddr.Port)
	if err != nil {
		log.Printf("invalid proxy server address: %v", err)
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

	// Register an output callback to write packets output from lwip stack to tun
	// device, output function should be set before input any packets.
	core.RegisterOutputFn(func(data []byte) (int, error) {
		// lwip -> tun
		return tunDev.Write(data)
	})

	stopReplyChannel = make(chan bool)
	stopSignalChannel = createDataPipeWorker()

	log.Printf("Running tun2socks")

	return 0
}
