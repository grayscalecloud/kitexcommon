package monitor

import "context"

type CtxCallback func(ctx context.Context)

type CtxErrCallback func(ctx context.Context) error
