package monitor

import (
	"github.com/grayscalecloud/kitexcommon/hdmodel"
)

func InitMonitor(serviceName string, cfg *hdmodel.Monitor) CtxCallback {
	if cfg.Enable {
		return initMetric(serviceName, cfg)
	}
	return nil
}
