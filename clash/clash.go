package clash

import (
	"os"
	"path/filepath"

	"github.com/Dreamacro/clash/config"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub/executor"

	log "github.com/sirupsen/logrus"
)

var (
	runningFlag bool
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

	if err := config.Init(C.Path.HomeDir()); err != nil {
		log.Fatalf("Initial configuration directory error: %s", err.Error())
	}

	cfg, err := executor.Parse()
	if err != nil {
		return
	}
	executor.ApplyConfig(cfg, true)
	runningFlag = true
	return
}

func IsRunning() bool {
	return runningFlag
}

func Stop() {
	//cancel()
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
	runningFlag = false
}
