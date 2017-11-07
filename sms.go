package sms

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/pborman/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var (
	phoneNumbersIsTooLong   = errors.New("PhoneNumbers is too long")
	phoneNumbersIsRequired  = errors.New("PhoneNumbers is required")
	accessKeyIdIsRequired   = errors.New("AccessKeyId is required")
	signNameIsRequired      = errors.New("SignName is required")
	templateCodeIsRequired  = errors.New("TemplateCode is required")
	templateParamIsRequired = errors.New("TemplateParam is required")
)

type Param struct {
	//system parameters
	AccessKeyId      string
	Timestamp        string
	Format           string
	SignatureMethod  string
	SignatureVersion string
	SignatureNonce   string
	Signature        string

	//business parameters
	Action        string
	Version       string
	RegionId      string
	PhoneNumbers  string
	SignName      string
	TemplateCode  string
	TemplateParam string
	OutId         string
}

type Client struct {
	EndPoint   string
	AccessId   string
	AccessKey  string
	HttpClient *http.Client
	Param      *Param
	param      map[string]string
}

type ErrorMessage struct {
	HttpCode  int    `json:"-"`
	RequestId string `json:"RequestId,omitempty"`
	Message   string `json:"Message,omitempty"`
	BizId     string `json:"BizId,omitempty"`
	Code      string `json:"Code,omitempty"`
}

func specialUrlEncode(value string) string {
	rstValue := url.QueryEscape(value)
	rstValue = strings.Replace(rstValue, "+", "%20", -1)
	rstValue = strings.Replace(rstValue, "*", "%2A", -1)
	rstValue = strings.Replace(rstValue, "%7E", "~", -1)

	return rstValue
}

func generateQueryStringAndSignature(params url.Values, accessKeySecret string) (string, string) {
	keys := make([]string, 0)
	for key, _ := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var buffer bytes.Buffer
	for _, key := range keys {
		buffer.WriteString("&" + specialUrlEncode(key) + "=" + specialUrlEncode(params.Get(key)))
	}

	sortQueryString := strings.TrimPrefix(buffer.String(), "&")
	stringToSign := "GET&" + specialUrlEncode("/") + "&" + specialUrlEncode(sortQueryString)
	sign := sign(accessKeySecret+"&", stringToSign)
	signature := specialUrlEncode(sign)

	return "&" + sortQueryString, signature
}

func sign(accessKeySecret, sortQueryString string) string {
	h := hmac.New(sha1.New, []byte(accessKeySecret))
	h.Write([]byte(sortQueryString))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (p *Param) ParamsIsValid() error {
	if len(p.PhoneNumbers) == 0 {
		return phoneNumbersIsRequired
	}
	if len(strings.Split(p.PhoneNumbers, ",")) > 27 {
		return phoneNumbersIsTooLong
	}
	if len(p.AccessKeyId) == 0 {
		return accessKeyIdIsRequired
	}
	if len(p.SignName) == 0 {
		return signNameIsRequired
	}
	if len(p.TemplateCode) == 0 {
		return templateCodeIsRequired
	}
	if len(p.TemplateParam) == 0 {
		return templateParamIsRequired
	}

	return nil
}

func (p *Param) BuildSmsRequestEndpoint(accessKeySecret, smsURL string) (string, error) {
	var err error
	if err = p.ParamsIsValid(); err != nil {
		return "", err
	}
	params := url.Values{}
	params.Set("AccessKeyId", p.AccessKeyId)
	params.Set("Timestamp", p.Timestamp)
	params.Set("SignatureMethod", p.SignatureMethod)
	params.Set("SignatureVersion", p.SignatureVersion)
	params.Set("SignatureNonce", p.SignatureNonce)
	params.Set("Format", p.Format)
	params.Set("Action", p.Action)
	params.Set("Version", p.Version)
	params.Set("RegionId", p.RegionId)
	params.Set("PhoneNumbers", p.PhoneNumbers)
	params.Set("SignName", p.SignName)
	params.Set("TemplateParam", p.TemplateParam)
	params.Set("TemplateCode", p.TemplateCode)
	params.Set("OutId", p.OutId)
	// generate signature and sorted  query
	sortQueryString, signature := generateQueryStringAndSignature(params, accessKeySecret)

	return smsURL + "?Signature=" + signature + sortQueryString, nil
}

func New(accessId, accessKey string) (c *Client) {
	c = new(Client)
	c.Param = &Param{}
	c.HttpClient = &http.Client{
		Timeout: time.Second * 3,
	}
	c.EndPoint = "http://dysmsapi.aliyuncs.com/"
	c.AccessId = accessId
	c.AccessKey = accessKey
	c.Param.Action = "SendSms"
	c.Param.Format = "JSON"
	c.Param.Version = "2017-05-25"
	c.Param.AccessKeyId = accessId
	c.Param.SignatureMethod = "HMAC-SHA1"
	c.Param.Timestamp = time.Now().UTC().Format(time.RFC3339)
	c.Param.SignatureVersion = "1.0"
	c.Param.SignatureNonce = uuid.New()
	c.Param.RegionId = "cn-hangzhou"
	c.Param.OutId = "out id"

	return c
}

func (c *Client) Send(phones, signName, templateCode, paramStr string) (e *ErrorMessage, err error) {
	c.Param.PhoneNumbers = phones
	c.Param.SignName = signName
	c.Param.TemplateCode = templateCode
	c.Param.TemplateParam = paramStr
	var endpoint string
	e = &ErrorMessage{}
	if endpoint, err = c.Param.BuildSmsRequestEndpoint(c.AccessKey, c.EndPoint); err != nil {
		return
	}
	request, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response, err := c.HttpClient.Do(request)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	defer response.Body.Close()

	err = json.Unmarshal(body, &e)

	return
}
