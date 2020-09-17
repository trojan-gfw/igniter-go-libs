package tun2socks

import (
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/eycorsican/go-tun2socks/common/dns/cache"
	"github.com/eycorsican/go-tun2socks/common/dns/fakedns"
	"github.com/eycorsican/go-tun2socks/common/log"
	"github.com/eycorsican/go-tun2socks/common/log/simpleandroidlog" // Register a simple android logger.
	"github.com/eycorsican/go-tun2socks/component/pool"
	"github.com/eycorsican/go-tun2socks/component/runner"
	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"

	"github.com/songgao/water"
)

type Tun2socksStartOptions struct {
	TunFd        int
	Socks5Server string
	FakeIPRange  string
	MTU          int
	EnableIPv6   bool
}

var (
	lwipWriter          io.Writer
	lwipStack           core.LWIPStack
	mtuUsed             int
	lwipTUNDataPipeTask *runner.Task
	tunDev              *water.Interface
)

// Stop stop it
func Stop() {
	log.Infof("enter stop")
	log.Infof("begin close tun")
	err := tunDev.Close()
	if err != nil {
		log.Infof("close tun: %v", err)
	}
	if lwipTUNDataPipeTask.Running() {
		log.Infof("send stop lwipTUNDataPipeTask sig")
		lwipTUNDataPipeTask.Stop()
		log.Infof("lwipTUNDataPipeTask stop sig sent")
		<-lwipTUNDataPipeTask.StopChan()
	} else {
		log.Infof("lwipTUNDataPipeTask already stopped")
	}

	log.Infof("begin close lwipStack")
	lwipStack.Close(core.DELAY)
}

// hack to receive tunfd
func openTunDevice(tunFd int) (*water.Interface, error) {
	file := os.NewFile(uintptr(tunFd), "tun") // dummy file path name since we already got the fd
	tunDev = &water.Interface{
		ReadWriteCloser: file,
	}
	return tunDev, nil
}

// Start sets up lwIP stack, starts a Tun2socks instance
func Start(opt *Tun2socksStartOptions) int {

	mtuUsed = opt.MTU
	var err error
	tunDev, err = openTunDevice(opt.TunFd)
	if err != nil {
		log.Fatalf("failed to open tun device: %v", err)
	}
	// handle previous lwIP stack
	if lwipStack != nil {
		log.Infof("begin close previous lwipStack")
		lwipStack.Close(core.INSTANT)
	} else {
		log.Infof("do NOT have to close previous lwipStack")
	}

	// Setup the lwIP stack.
	lwipStack = core.NewLWIPStack(opt.EnableIPv6)
	lwipWriter = lwipStack.(io.Writer)

	// Register tun2socks connection handlers.
	proxyAddr, err := net.ResolveTCPAddr("tcp", opt.Socks5Server)
	proxyHost := proxyAddr.IP.String()
	proxyPort := uint16(proxyAddr.Port)
	if err != nil {
		log.Infof("invalid proxy server address: %v", err)
		return -1
	}
	cacheDNS := cache.NewSimpleDnsCache()
	if opt.FakeIPRange != "" {
		_, ipnet, err := net.ParseCIDR(opt.FakeIPRange)
		if err != nil {
			log.Fatalf("failed to parse fake ip range %v", opt.FakeIPRange)
		}
		fakeDNS := fakedns.NewFakeDNS(ipnet, 3000)
		core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, proxyPort, fakeDNS))
		core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, proxyPort, 30*time.Second, cacheDNS, fakeDNS))
	} else {
		core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, proxyPort, nil))
		core.RegisterUDPConnHandler(socks.NewUDPHandler(proxyHost, proxyPort, 30*time.Second, cacheDNS, nil))
	}

	// Register an output callback to write packets output from lwip stack to tun
	// device, output function should be set before input any packets.
	core.RegisterOutputFn(func(data []byte) (int, error) {
		// lwip -> tun
		return tunDev.Write(data)
	})

	if lwipTUNDataPipeTask != nil && lwipTUNDataPipeTask.Running() {
		log.Infof("stop previous lwipTUNDataPipeTask sig")
		lwipTUNDataPipeTask.Stop()
		log.Infof("previous lwipTUNDataPipeTask stop sig sent")
		<-lwipTUNDataPipeTask.StopChan()
	} else {
		log.Infof("previous lwipTUNDataPipeTask already stopped or never being started")
	}

	lwipTUNDataPipeTask = runner.Go(func(shouldStop runner.S) error {
		// do setup
		// defer func(){
		//   // do teardown
		// }
		zeroErr := errors.New("no error")
		maxErrorTimes := 20
		for {
			// do some work here

			// tun -> lwip
			err := func() error {
				buf := pool.NewBytes(pool.BufSize)
				defer pool.FreeBytes(buf)
				_, err := io.CopyBuffer(lwipWriter, tunDev, buf)
				return err
			}()
			if err != nil {
				maxErrorTimes--
				log.Infof("copying data failed: %v", err)
			}
			if shouldStop() {
				log.Infof("got DataPipe stop signal")
				break
			}
			if maxErrorTimes <= 0 {
				log.Infof("lwipTUNDataPipeTask returns due to exceeded error times")
				return err
			}
		}
		log.Infof("exit DataPipe loop")
		return zeroErr // any errors?
	})

	log.Infof("Running tun2socks")

	return 0
}

// SetLoglevel set tun2socks log level
// possible input: debug/info/warn/error/none
func SetLoglevel(logLevel string) {
	// Set log level.
	switch strings.ToLower(logLevel) {
	case "debug":
		log.SetLevel(log.DEBUG)
	case "info":
		log.SetLevel(log.INFO)
	case "warn":
		log.SetLevel(log.WARN)
	case "error":
		log.SetLevel(log.ERROR)
	case "none":
		log.SetLevel(log.NONE)
	default:
		panic("unsupport logging level")
	}
	logger := simpleandroidlog.GetLogger()
	log.Infof("LogLevel: %v", logger.GetLevel())
}
