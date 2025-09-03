package errno

import "fmt"

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}
func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

// NewErrNoWithParams return ErrNo
func NewErrNoWithParams(code int64, msg string) *ErrNo {
	return &ErrNo{
		ErrCode: code,
		ErrMsg:  msg,
	}
}
