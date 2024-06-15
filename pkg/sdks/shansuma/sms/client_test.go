package sms

import (
	"testing"

	"github.com/spf13/viper"
)

func TestClient_Execute(t *testing.T) {
	params := make(map[string]interface{})
	params["code"] = 123456

	client := NewSmsProvider(&viper.Viper{})

	request := NewRequest()
	request.SetMethod("sms.message.send")
	request.SetBizContent(TemplateMessage{
		Mobile:     []string{"18374775592"},
		Type:       0,
		Sign:       "案例笔记",
		TemplateId: "ST_2022060400000001",
		SendTime:   "",
		Params:     params,
	})

	err := client.Execute(request)
	t.Log("err : ", err)
}
