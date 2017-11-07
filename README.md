# aliyun-sms-sdk-golang
基于HTTP协议的阿里云短信服务golang版本实现

## About

消息服务同时具备发送短信的能力，支持快速发送短信验证码、短信通知、推广短信。完美支撑双11期间的2亿用户发送6亿条短信。三网合一专属通道，与工信部携号转网平台实时互联。电信级运维保证，实时监控自动切换，到达率高达99%。

> 阿里云消息服务（Message Service，原MQS）是阿里云商用的消息中间件服务。与传统的消息中间件不同，消息服务一开始就是基于阿里云自主研发的飞天分布式系统来设计和实现，具有大规模，高可靠、高并发访问和超强消息堆积能力的特点。消息服务API采用HTTP RESTful标准，接入方便，跨网络能力强；已全面接入资源访问控制服务（RAM）、专有网络（VPC），支持各种安全访问控制；接入云监控，提供完善的监控及报警机制。消息服务提供丰富的SDK、解决方案、最佳实践和7x24小时的技术支持，帮助应用开发者在应用组件之间自由地传递数据和构建松耦合、分布式、高可用系统。

## Install

```go
go get -u github.com/JakeXu/aliyun-sms-sdk-golang
```

## Usage

[使用帮助](doc/tips.md)

```go

package main

import (
	"github.com/JakeXu/aliyun-sms-sdk-golang"
	"fmt"
)

// modify it to yours
const (
	ACCESSID  = "your_accessid"
	ACCESSKEY = "your_accesskey"
)

func main() {
	c := sms.New(ACCESSID, ACCESSKEY)
	e, err := c.Send("1380000****", "阿里云短信测试专用", "SMS_71390007", `{"code":"123456"}`)
	if err != nil {
        fmt.Println("send sms failed", err, e.Message)
    } else {
        fmt.Println("send sms succeed", e.RequestId)
        // business logic operation
    }
}

```

## Links

* [Short Message Service，SMS(短信服务)](https://www.aliyun.com/product/mns?spm=5176.8195934.765261.239.1c35a6f94ZT1OV)
* [API使用手册](https://help.aliyun.com/document_detail/56189.html?spm=5176.8195934.507901.13.0ywNYu)