package monitor

import (
	"github.com/grayscalecloud/kitexcommon/hdmodel"
)

func InitMonitor(serviceName string, cfg *hdmodel.Monitor) CtxCallback {
	if cfg.Enable {
		InitTracing(serviceName)
		return initMetric(serviceName, cfg)
	}
	return nil
}
