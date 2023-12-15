//go:build linux
// +build linux

package log

import (
	"fmt"
	"github.com/JoeyZeYi/source/log/zap"
	"os"
	"path/filepath"
)

func InitLog() {
	level := zap.NewAtomicLevel()
	InitAppLog(LogPath(fmt.Sprintf("/data/project/log/%s.log", filepath.Base(os.Args[0]))), TestEnv(false), SetLevel(level))
}
