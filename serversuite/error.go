package serversuite

import (
	"context"
	"errors"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/remote"
)

// convert errors that can be serialized
func ServerErrorHandler(ctx context.Context, err error) error {
	// if you want get other rpc info, you can get rpcinfo first, like `ri := rpcinfo.GetRPCInfo(ctx)`
	// for example, get remote address: `remoteAddr := rpcinfo.GetRPCInfo(ctx).From().Address()`

	if errors.Is(err, kerrors.ErrBiz) {
		err = errors.Unwrap(err)
	}
	if errCode, ok := GetErrorCode(err); ok {
		// for Thrift、KitexProtobuf
		return remote.NewTransError(errCode, err)
	}
	return err
}

func GetErrorCode(err error) (int32, bool) {
	return 0, false
}
