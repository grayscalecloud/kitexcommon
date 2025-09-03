package clientsuite

import (
	"context"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/remote"
)

func ClientErrorHandler(ctx context.Context, err error) error {
	// if you want get other rpc info, you can get rpcinfo first, like `ri := rpcinfo.GetRPCInfo(ctx)`
	// for example, get remote address: `remoteAddr := rpcinfo.GetRPCInfo(ctx).To().Address()`

	// for thrift、KitexProtobuf
	if e, ok := err.(*remote.TransError); ok {
		// TypeID is error code
		return buildYourError(e.TypeID(), e)
	}

	return kerrors.ErrRemoteOrNetwork.WithCause(err)
}

func buildYourError(id int32, e *remote.TransError) error {
	return nil
}
