package sms

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"strings"
	"testing"
)

const (
	accessKeyId      string = "testId"
	accessKeySecret  string = "testSecret"
	format           string = "JSON"
	action           string = "SendSms"
	version          string = "2017-05-25"
	regionId         string = "cn-hangzhou"
	endPoint         string = "http://dysmsapi.aliyuncs.com/"
	signatureMethod  string = "HMAC-SHA1"
	signatureVersion string = "1.0"
	signName         string = "阿里云短信测试专用"
	templateCode     string = "SMS_71390007"
	templateParam    string = "{\"code\":\"123456\"}"
	phoneNumbers     string = "15300000001"
)

func Test_specialUrlEncode(t *testing.T) {
	assert := assert.New(t)
	encoding := specialUrlEncode("+*%7E")
	assert.NotContains(encoding, "+")
	assert.NotContains(encoding, "*")
	assert.NotContains(encoding, "%7E")
}

func Test_generateQueryStringAndSignature(t *testing.T) {
	assert := assert.New(t)
	params := url.Values{}
	params.Set("AccessKeyId", accessKeyId)
	params.Set("Format", format)

	sortQueryString, signature := generateQueryStringAndSignature(params, accessKeySecret)

	assert.Contains(sortQueryString, "AccessKeyId")
	assert.Contains(sortQueryString, "Format")
	assert.Contains(signature, "%3D")
}

func Test_New(t *testing.T) {
	assert := assert.New(t)
	c := New(accessKeyId, accessKeySecret)

	assert.Equal(action, action)
	assert.Equal(version, version)
	assert.Equal(regionId, regionId)
	assert.Equal(endPoint, endPoint)
	assert.Equal(signatureMethod, signatureMethod)
	assert.Equal(signatureVersion, signatureVersion)
	phones := make([]string, 28)
	for k, _ := range phones {
		phones[k] = phoneNumbers
	}
	c.Param.PhoneNumbers = strings.Join(phones, ",")
	assert.Equal(c.Param.ParamsIsValid(), phoneNumbersIsTooLong)

	c.Param.PhoneNumbers = ""
	assert.Equal(c.Param.ParamsIsValid(), phoneNumbersIsRequired)

	c.Param.AccessKeyId = ""
	c.Param.PhoneNumbers = phoneNumbers
	assert.Equal(c.Param.ParamsIsValid(), accessKeyIdIsRequired)

	c.Param.AccessKeyId = accessKeyId
	assert.Equal(c.Param.ParamsIsValid(), signNameIsRequired)

	c.Param.SignName = signName
	assert.Equal(c.Param.ParamsIsValid(), templateCodeIsRequired)

	c.Param.TemplateCode = templateCode
	assert.Equal(c.Param.ParamsIsValid(), templateParamIsRequired)

	c.Param.TemplateParam = templateParam
	assert.Nil(c.Param.ParamsIsValid())
}

func Test_BuildSmsRequestEndpoint(t *testing.T) {
	assert := assert.New(t)
	c := New(accessKeyId, accessKeySecret)

	c.Param.PhoneNumbers = phoneNumbers
	c.Param.TemplateCode = templateCode
	c.Param.TemplateParam = templateParam
	endpoint, err := c.Param.BuildSmsRequestEndpoint(c.AccessKey, c.EndPoint)
	assert.NotNil(err)

	c.Param.SignName = signName
	endpoint, err = c.Param.BuildSmsRequestEndpoint(c.AccessKey, c.EndPoint)
	assert.Nil(err)
	assert.Contains(endpoint, "Signature")
	assert.Contains(endpoint, "AccessKeyId")
	assert.Contains(endpoint, "Timestamp")
	assert.Contains(endpoint, "SignatureMethod")
	assert.Contains(endpoint, "SignatureVersion")
	assert.Contains(endpoint, "SignatureNonce")
	assert.Contains(endpoint, "Format")
	assert.Contains(endpoint, "Action")
	assert.Contains(endpoint, "RegionId")
	assert.Contains(endpoint, "PhoneNumbers")
	assert.Contains(endpoint, "SignName")
	assert.Contains(endpoint, "TemplateParam")
	assert.Contains(endpoint, "TemplateCode")
	assert.Contains(endpoint, "OutId")
}

func Test_Send(t *testing.T) {
	assert := assert.New(t)

	c := New(accessKeyId, accessKeySecret)
	msg, err := c.Send("", signName, templateCode, templateParam)
	assert.NotNil(err)

	msg, err = c.Send(phoneNumbers, signName, templateCode, templateParam)
	assert.NotEmpty(msg.Message)
	assert.Nil(err)
}
