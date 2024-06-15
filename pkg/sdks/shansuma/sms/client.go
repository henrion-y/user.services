package sms

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Client struct {
	appId      string
	secretKey  string
	templateId string
	sign       string
	signType   string
	version    string
	env        string
}

func NewSmsProvider(config *viper.Viper) *Client {
	appId := config.GetString("sms.AppId")
	secretKey := config.GetString("sms.SecretKey")
	templateId := config.GetString("sms.TemplateId")
	sign := config.GetString("sms.Sign")
	env := config.GetString("sms.Env")
	c := &Client{
		env:        env,
		sign:       sign,
		appId:      appId,
		secretKey:  secretKey,
		templateId: templateId,
		version:    "1.0",
		signType:   "md5",
	}
	return c
}

func (c *Client) SendCode(ctx context.Context, phone string, code string) error {
	if c.env == "debug" {
		return nil
	}
	params := make(map[string]interface{})
	params["code"] = code
	request := NewRequest()
	request.SetMethod("sms.message.send")
	request.SetBizContent(TemplateMessage{
		Mobile:     []string{phone},
		Type:       0,
		Sign:       c.sign,
		TemplateId: c.templateId,
		SendTime:   "",
		Params:     params,
	})
	return c.Execute(request)
}

func (c *Client) CreateSignature(data map[string]string) string {
	keys := make([]string, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	list := make([]string, 0)
	for _, k := range keys {
		list = append(list, fmt.Sprintf("%s=%s", k, data[k]))
	}
	list = append(list, fmt.Sprintf("%s=%s", "key", c.secretKey))

	str := strings.Join(list, "&")
	str = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(str))))
	return str
}

func (c *Client) Execute(request *Request) error {
	post := make(map[string]string)
	post["app_id"] = c.appId
	post["method"] = request.GetMethod()
	post["version"] = c.version
	post["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	post["sign_type"] = c.signType
	post["biz_content"] = request.GetBizContent()
	post["sign"] = c.CreateSignature(post)
	data := url.Values{}
	for name, value := range post {
		data.Set(name, value)
	}
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://api.shansuma.com/gateway.do", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla 5.0 GO-SMS-SDK")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	result := &Response{}
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, result)
	if err != nil {
		return err
	}
	if result.Result.Code == 0 {
		return nil
	}
	return errors.New(result.Result.Message)
}
