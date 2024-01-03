package log

import (
	"github.com/JoeyZeYi/source/log/zap"
	"os"
	"path"
)

const logApiPath = "/log"

type appZapLogConf struct {
	testEnv         bool
	processName     string
	withPid         bool
	logApiPath      string
	listenAddr      string
	HostName        string
	ElkTemplateName string //区分不同业务
	logPath         string
	Level           zap.AtomicLevel
	maxSize         int  //文件存储大小
	maxAge          int  //文件保留时间 day
	gzip            bool //是否启用Gzip压缩
}

var defaultLogOptions = appZapLogConf{
	testEnv:         true,
	processName:     path.Base(os.Args[0]),
	withPid:         true,
	logApiPath:      logApiPath,
	listenAddr:      "127.0.0.1:0",
	ElkTemplateName: path.Base(os.Args[0]),
}

type IAppZapOption interface {
	apply(*appZapLogConf)
}

type appZapOptionFunc func(*appZapLogConf)

func (app appZapOptionFunc) apply(option *appZapLogConf) {
	app(option)
}

// ListenAddr ListenAddr设置logserver的http端口,用来管理日志级别。
// 默认监听127.0.0.1下的随机端口
func ListenAddr(addr string) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.listenAddr = addr
	})
}

// LogApiPath LogApiPath设置logserver的api名字。
// 默认为 /log。
func LogApiPath(apiPath string) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.logApiPath = apiPath
	})
}

func SetLevel(level zap.AtomicLevel) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.Level = level
	})
}

// LogPath 设置日志路径
func LogPath(path string) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.logPath = path
	})
}

// WithPid WithPid设置日志输出中是否加入pid的项。
// 默认为true。
func WithPid(yes bool) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.withPid = yes
	})
}

// ProcessName ProcessName设置输出的进程名字。
// 默认去当前执行文件的名字。
func ProcessName(name string) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.processName = name
	})
}

// TestEnv TestEnv设置是否测试环境。
// 默认为true,测试环境。
func TestEnv(yes bool) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.testEnv = yes
	})
}

// HostName HostName设置日志机器的ip地址,方便定位。
// 默认不输出。
func HostName(hostname string) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.HostName = hostname
	})
}

// MaxSize 文件存储大小 单位是M 默认500M
func MaxSize(maxSize int) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.maxSize = maxSize
	})
}

// MaxAge 文件保留时间 单位是天 默认一直保留
func MaxAge(maxAge int) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.maxAge = maxAge
	})
}

// Gzip 是否启用Gzip压缩
func Gzip(gzip bool) IAppZapOption {
	return appZapOptionFunc(func(option *appZapLogConf) {
		option.gzip = gzip
	})
}
