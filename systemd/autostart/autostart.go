//go:build linux
// +build linux

package autostart

import (
	"github.com/JoeyZeYi/source/log"
	"github.com/JoeyZeYi/source/log/zap"
	"github.com/JoeyZeYi/source/systemd"
	"os"
	"path/filepath"
)

func AutostartInit() {
	if err := systemd.AutoStart(filepath.Base(os.Args[0])); err != nil {
		log.Fatal("systemd.AutoStart", zap.Error(err))
	}
}
