package sms

import (
	"fmt"
	"github.com/JakeXu/aliyun-sms-sdk-golang"
)

// 验证码短信接口实例化, 方便使用
func CaptchaSMS() *sms.Client {
	return sms.New("testId", "testSecret")
}

func testSMS() {
	e, err := CaptchaSMS().Send("15300000001", "阿里云短信测试专用", "SMS_71390007", "{\"code\":\"123456\"}")
	if err != nil {
		fmt.Println("send sms failed", err, e.Message)
	} else {
		fmt.Println("send sms succeed", e.RequestId)
		// business logic operation
	}
}
