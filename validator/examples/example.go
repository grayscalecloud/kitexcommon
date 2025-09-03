package main

import (
	"fmt"

	"github.com/grayscalecloud/kitexcommon/validator"
)

// User 用户模型
type User struct {
	Username  string
	Email     string
	Password  string
	Age       int
	Phone     string
	Interests []string
}

func main() {
	// 创建验证器
	v := validator.NewValidator()

	// 添加验证规则
	v.AddRule("Username", &validator.Required{Message: "用户名不能为空"})
	v.AddRule("Username", &validator.MinLength{Length: 3, Message: "用户名长度不能小于3个字符"})
	v.AddRule("Username", &validator.MaxLength{Length: 20, Message: "用户名长度不能超过20个字符"})

	v.AddRule("Email", &validator.Required{Message: "邮箱不能为空"})
	v.AddRule("Email", &validator.Email{Message: "邮箱格式不正确"})

	v.AddRule("Password", &validator.Required{Message: "密码不能为空"})
	v.AddRule("Password", &validator.MinLength{Length: 8, Message: "密码长度不能小于8个字符"})
	v.AddRule("Password", &validator.Pattern{
		Regex:   "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d).*$",
		Message: "密码必须包含大小写字母和数字",
	})

	v.AddRule("Age", &validator.CustomRule{
		Func: func(value interface{}) (bool, string) {
			age, ok := value.(int)
			if !ok {
				return false, "年龄必须是整数"
			}
			if age < 18 {
				return false, "年龄必须大于或等于18岁"
			}
			if age > 120 {
				return false, "年龄不合理"
			}
			return true, ""
		},
	})

	v.AddRule("Phone", &validator.Phone{Message: "手机号格式不正确"})

	v.AddRule("Interests", &validator.Required{Message: "兴趣爱好不能为空"})
	v.AddRule("Interests", &validator.MinLength{Length: 1, Message: "至少需要一个兴趣爱好"})

	// 验证有效数据
	validUser := User{
		Username:  "johndoe",
		Email:     "john.doe@example.com",
		Password:  "Password123",
		Age:       25,
		Phone:     "13800138000",
		Interests: []string{"reading", "coding"},
	}

	valid, errors := v.Validate(validUser)
	fmt.Println("验证有效用户:")
	fmt.Printf("是否有效: %v\n", valid)
	if !valid {
		fmt.Printf("错误信息: %s\n", validator.FormatErrors(errors))
	}
	fmt.Println()

	// 验证无效数据
	invalidUser := User{
		Username:  "jo",
		Email:     "not-an-email",
		Password:  "weak",
		Age:       15,
		Phone:     "12345",
		Interests: []string{},
	}

	valid, errors = v.Validate(invalidUser)
	fmt.Println("验证无效用户:")
	fmt.Printf("是否有效: %v\n", valid)
	if !valid {
		fmt.Printf("错误信息: %s\n", validator.FormatErrors(errors))
	}
}
