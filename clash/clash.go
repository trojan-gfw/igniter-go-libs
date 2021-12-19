package clash

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Dreamacro/clash/config"
	C "github.com/Dreamacro/clash/constant"
	_ "github.com/Dreamacro/clash/hub"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/tunnel/statistic"

	log "github.com/sirupsen/logrus"
)

type ClashStartOptions struct {
	// HomeDir Clash config home directory
	HomeDir string
	// SocksListener Clash listener address and port
	SocksListener string
	// TrojanProxyServer Trojan proxy listening address and port
	TrojanProxyServer string
	// TrojanProxyServerUdpEnabled Whether UDP is enabled for Trojan Server
	TrojanProxyServerUdpEnabled bool
}

func Start(opt *ClashStartOptions) {
	homedir := opt.HomeDir
	if homedir != "" {
		if !filepath.IsAbs(homedir) {
			currentDir, _ := os.Getwd()
			homedir = filepath.Join(currentDir, homedir)
		}
		C.SetHomeDir(homedir)
	}

	configFile := filepath.Join(C.Path.HomeDir(), "config.yaml")
	C.SetConfig(configFile)

	if err := config.Init(C.Path.HomeDir()); err != nil {
		log.Fatalf("Initial configuration directory error: %s", err.Error())
	}

	ApplyRawConfig(opt)
	return
}

func Stop() {
	snapshot := statistic.DefaultManager.Snapshot()
	for _, c := range snapshot.Connections {
		err := c.Close()
		if err != nil {
			log.Warnf("Clash Stop(): close conn err %v", err)
		}
	}

	opt := &ClashStartOptions{
		SocksListener:               "127.0.0.1:0",
		TrojanProxyServer:           "127.0.0.1:0",
		TrojanProxyServerUdpEnabled: true,
	}
	ApplyRawConfig(opt)
}

func ApplyRawConfig(opt *ClashStartOptions) {

	// handle user input
	socksListenerHost, socksListenerPort, err := net.SplitHostPort(opt.SocksListener)
	if err != nil {
		log.Fatalf("SplitHostPort err: %v (%v)", err, opt.SocksListener)
	}
	if len(socksListenerHost) <= 0 {
		log.Fatalf("SplitHostPort host is empty: %v", socksListenerHost)
	}
	trojanProxyServerHost, trojanProxyServerPort, err := net.SplitHostPort(opt.TrojanProxyServer)
	if err != nil {
		log.Fatalf("SplitHostPort err: %v (%v)", err, opt.TrojanProxyServer)
	}
	if len(trojanProxyServerHost) <= 0 {
		log.Fatalf("SplitHostPort host is empty: %v", trojanProxyServerHost)
	}

	rawConfigBytes, err := readConfig(C.Path.Config())
	if err != nil {
		log.Fatalf("fail to read Clash config file: %v", err)
	}
	rawCfg, err := config.UnmarshalRawConfig(rawConfigBytes)
	if err != nil {
		log.Fatalf("UnmarshalRawConfig: %v", err)
	}

	port, err := strconv.Atoi(socksListenerPort)
	if err != nil {
		log.Fatalf("fail to convert socksListenerPort %v", socksListenerPort)
	}
	if len(rawCfg.Proxy) <= 0 {
		log.Fatalf("should at least add one upstream proxy server")
	}

	rawCfg.AllowLan = true // whether we really use this feature is determined by BindAddress
	rawCfg.SocksPort = port
	rawCfg.BindAddress = socksListenerHost //default is *
	firstProxyServerMap := rawCfg.Proxy[0]
	//proxies:
	//  - { name: "trojan", type: socks5, server: "127.0.0.1", port: 1081, udp: true}
	if firstProxyServerMap["type"] == "socks5" && firstProxyServerMap["name"] == "trojan" {
		firstProxyServerMap["server"] = trojanProxyServerHost
		firstProxyServerMap["port"] = trojanProxyServerPort
		firstProxyServerMap["udp"] = opt.TrojanProxyServerUdpEnabled
	} else {
		log.Fatalf("fail to find trojan proxy entry in Clash config")
	}

	cfg, err := config.ParseRawConfig(rawCfg)
	if err != nil {
		log.Fatalf("ParseRawConfig: %v", err)
	}

	executor.ApplyConfig(cfg, true)
}

func readConfig(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("configuration file %s is empty", path)
	}

	return data, err
}
