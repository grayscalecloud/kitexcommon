package monitor

import "github.com/grayscalecloud/kitexcommon/model"

func Init(serviceName string, cfg *model.Monitor) CtxCallback {
	if cfg.Enabled {
		return initMetric(serviceName, cfg)
	}
	return nil
}
