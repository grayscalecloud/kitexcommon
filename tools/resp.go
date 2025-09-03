package tools

import (
	"errors"

	"github.com/bytedance/sonic"
	"github.com/grayscalecloud/kitexcommon/base"
	"github.com/grayscalecloud/kitexcommon/consts/errno"
)

func BuildBaseResp(err error) *base.BaseResponse {
	if err == nil {
		return baseResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(&e)
	}
	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(&s)
}

func baseResp(err *errno.ErrNo) *base.BaseResponse {
	return &base.BaseResponse{
		Code: err.ErrCode,
		Msg:  err.ErrMsg,
	}
}

func BuildBaseRespSuccess(msg string) *base.BaseResponse {
	return &base.BaseResponse{
		Code: errno.Success.ErrCode,
		Msg:  msg,
	}
}

func BuildBaseRespSuccessNoParams() *base.BaseResponse {
	return &base.BaseResponse{
		Code: errno.Success.ErrCode,
		Msg:  "Success",
	}
}

func BuildBaseRespFailNoParams() *base.BaseResponse {
	return &base.BaseResponse{
		Code: errno.BadRequest.ErrCode,
		Msg:  "Fail",
	}
}

// ConventKitexToHertz 转换kitex返回的response为hertz的response
func ConventKitexToHertz[T any](kitexResp interface{}) (hertzResp T, err error) {
	marshal, err := sonic.Marshal(kitexResp)
	if err != nil {
		return hertzResp, err
	}
	err = sonic.Unmarshal(marshal, &hertzResp)
	if err != nil {
		return hertzResp, err
	}
	return hertzResp, nil
}
