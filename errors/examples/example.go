package main

import (
	"fmt"
	"os"

	"github.com/grayscalecloud/kitexcommon/errors"
)

// 定义应用程序错误码
const (
	ErrNotFound      = 404
	ErrUnauthorized  = 401
	ErrInternalError = 500
)

func main() {
	// 示例1: 创建新错误
	err1 := errors.NewError(ErrNotFound, "用户不存在")
	fmt.Println("示例1 - 创建新错误:")
	fmt.Println(err1.Error())
	fmt.Println(err1.FormatStack())
	fmt.Println()

	// 示例2: 包装标准错误
	_, err2 := os.Open("不存在的文件.txt")
	wrappedErr := errors.Wrap(err2, ErrNotFound, "无法打开文件")
	fmt.Println("示例2 - 包装标准错误:")
	fmt.Println(wrappedErr.Error())
	fmt.Println()

	// 示例3: 错误类型检查
	if errors.Is(wrappedErr, os.ErrNotExist) {
		fmt.Println("示例3 - 错误类型检查: 文件确实不存在")
	}
	fmt.Println()

	// 示例4: 错误链
	err4 := someFunction()
	fmt.Println("示例4 - 错误链:")
	fmt.Println(err4.Error())
	fmt.Println()

	// 示例5: 获取错误信息
	//if e, ok := err4.(*errors.Error); ok {
	//	fmt.Printf("示例5 - 获取错误信息:\n")
	//	fmt.Printf("错误码: %d\n", e.GetCode())
	//	fmt.Printf("错误消息: %s\n", e.GetMessage())
	//	fmt.Println(e.FormatStack())
	//}
}

// 模拟数据库操作
func queryDatabase() error {
	return errors.NewError(ErrInternalError, "数据库连接失败")
}

// 模拟业务逻辑
func businessLogic() error {
	err := queryDatabase()
	if err != nil {
		return errors.WrapWithMessage(err, "执行业务逻辑时出错")
	}
	return nil
}

// 模拟API处理
func someFunction() *errors.Error {
	err := businessLogic()
	if err != nil {
		// 类型断言，如果不是*errors.Error类型，则创建一个新的
		if e, ok := err.(*errors.Error); ok {
			return errors.Wrap(e, ErrInternalError, "API处理失败")
		}
		return errors.Wrap(err, ErrInternalError, "API处理失败")
	}
	return nil
}
