package sms_jatis

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	CHANNEL_REGULAR = "0"
	CHANNEL_ALERT   = "1"
	CHANNEL_OTP     = "2"
)

// Config config app for Jatis
type Config struct {
	BaseUrl        string
	UserId         string
	Password       string
	Sender         string
	Division       string
	UploadBy       string
	channel        int
	ConnectTimeout int
}

type Sender struct {
	Config Config
}

// end of config

type ResponseBody struct {
	MessageId string
	To        string
	Status    string
}

// end of response

// ReqMessage Request for client
type ReqMessage struct {
	PhoneNumber string
	Text        string
}

// end of request

func New(config Config) *Sender {
	return &Sender{
		Config: config,
	}
}

// SendSMS function to send message
func (s *Sender) SendSMS(ctx context.Context, request ReqMessage) (respBody ResponseBody, err error) {
	urlPath := fmt.Sprintf("%s", s.Config.BaseUrl)
	data := url.Values{}
	data.Set("userid", s.Config.UserId)
	data.Set("password", s.Config.Password)
	data.Set("msisdn", request.PhoneNumber)
	data.Set("message", request.Text)
	data.Set("sender", s.Config.Sender)
	data.Set("batchname", "otp")
	data.Set("division", s.Config.Division)
	data.Set("uploadby", s.Config.UploadBy)
	data.Set("channel", CHANNEL_OTP)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error(err)
		err = errors.New(JatisError[ErrorInternal])
		return
	}

	req.Header.Set("Content-Type", "x-www-form-urlencoded")

	timeout := s.Config.ConnectTimeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		err = errors.New(JatisError[ErrorInternal])
		log.Error(err)
	}

	if res.StatusCode >= 400 {
		log.Error(res)
		err = errors.New(JatisError[ErrorInternal])
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(res)
		err = errors.New(JatisError[ErrorInternal])
		return
	}

	respMap, err := url.ParseQuery(string(body))
	if err != nil {
		return
	}

	respBody.To = request.PhoneNumber
	respBody.Status = respMap["Status"][0]
	// ommit error checking, worst case will be zero
	respStatus, _ := strconv.Atoi(respBody.Status)
	if respStatus != StatusSuccess {
		// return human readable error message
		respBody.Status = JatisError[respStatus]
		err = errors.New(JatisError[respStatus])
		return
	}
	respBody.MessageId = respMap["MessageId"][0]

	return respBody, nil
}
