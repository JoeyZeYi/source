package log

import (
	"github.com/JoeyZeYi/source/log/zap"
	"net"
	"net/http"
)

var setLevelPath string

func logLevelHttpServer(config *appZapLogConf, level zap.AtomicLevel) {

	logServerMux := http.NewServeMux()
	// set log level
	logServerMux.Handle(config.logApiPath, level)
	listener, err := net.Listen("tcp", config.listenAddr)
	if err != nil {
		Fatal("failed Listen: ",
			zap.String("ipport", config.listenAddr),
			zap.Error(err),
		)
	} else {
		setLevelPath = "http://" + listener.Addr().String() + config.logApiPath
		Info("open log service success", zap.String("url", setLevelPath))
	}
	go func() {
		err = http.Serve(listener, logServerMux)
		if err != nil {
			Fatal("failed ListenAndServe: ",
				zap.String("ipport", net.JoinHostPort("127.0.0.1", "18001")),
				zap.Error(err),
			)
		}
	}()

}
