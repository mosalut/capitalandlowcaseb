// This file is auto-generated, don't edit it. Thanks.
package main

import (
//	"os"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	openapi	"github.com/alibabacloud-go/darabonba-openapi/client"
	"github.com/alibabacloud-go/tea/tea"
	"fmt"
)

type sms_T struct {
	PhoneNumbers string
	SignName string
	TemplateCode string
	TemplateParam string
	OutId string
}

type smsStorage_T struct {
	code string
	timestamp int64
}

var smsM = make(map[string]*smsStorage_T)

func (sms *sms_T) send() error {
	client, err := CreateClient(tea.String("LTAI5tPZEHrkFBbskXpp7gYm"), tea.String("jguZm5G02DtXqRUbeaAYlKMrs8dHpK"))
	if err != nil {
		return err
	}

	fmt.Println(sms.PhoneNumbers)

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers: tea.String(sms.PhoneNumbers),
		SignName: tea.String(sms.SignName),
		TemplateCode: tea.String(sms.TemplateCode),
		TemplateParam: tea.String(sms.TemplateParam),
		OutId: tea.String(sms.OutId),
	}
	// 复制代码运行请自行打印 API 的返回值
	resp, err := client.SendSms(sendSmsRequest)
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return err
}

func CreateClient (accessKeyId *string, accessKeySecret *string) (*dysmsapi20170525.Client, error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	result := &dysmsapi20170525.Client{}
	result, err := dysmsapi20170525.NewClient(config)
	return result, err
}
