package clash

import (
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/Dreamacro/clash/config"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/tunnel"

	log "github.com/sirupsen/logrus"
)

var (
	runningFlag atomic.Value
)

func Start(homedir string) {
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")

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

	cfg, err := executor.Parse()
	if err != nil {
		return
	}
	executor.ApplyConfig(cfg, true)
	runningFlag.Store(true)
	return
}

func IsRunning() bool {
	run := runningFlag.Load()
	return run.(bool)
}

func Stop() {
	// this is an unofficial feature of Clash, from: https://github.com/Dreamacro/clash/pull/341
	g := &config.General{
		Port:      0,
		SocksPort: 0,
	}
	cfg := &config.Config{
		General:      g,
		DNS:          &config.DNS{},
		Experimental: &config.Experimental{},
	}

	executor.ApplyConfig(cfg, true)

	snapshot := tunnel.DefaultManager.Snapshot()
	for _, c := range snapshot.Connections {
		c.Close()
	}

	runningFlag.Store(false)
}
