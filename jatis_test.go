package sms_jatis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// integration test
func TestSendSmsSuccess(t *testing.T) {
	config := Config{
		BaseUrl:        "https://sms-api.jatismobile.com/index.ashx",
		UserId:         "OperasionalFLIP",
		Password:       "OperasionalFLIP1120",
		Sender:         "Flip",
		Division:       "Divisi Operasional",
		UploadBy:       "adeputriawuser123",
		channel:        2,
		ProxyUrl:       "",
		ConnectTimeout: 15,
	}
	sender := New(config)
	reqMsg := ReqMessage{
		PhoneNumber: "089662233555",
		Text:        "tes",
	}

	res, err := sender.SendSMS(reqMsg)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}