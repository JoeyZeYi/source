//go:build windows
// +build windows

package log

import (
	"github.com/JoeyZeYi/source/log/zap"
)

func InitLog() {
	level := zap.NewAtomicLevel()
	InitAppLog(SetLevel(level))
}
