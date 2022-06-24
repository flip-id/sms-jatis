package sms_jatis

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const CURLOPT_PROXY = "10004"

// Config config app for Jatis
type Config struct {
	BaseUrl        string
	UserId         string
	Password       string
	Sender         string
	Division       string
	UploadBy       string
	channel        int
	ProxyUrl       string
	ConnectTimeout int
}

type Sender struct {
	Config Config
}

// end of config

// Message Response from Jatis
type Message struct {
	To     string
	Status string
}

type ResponseBody struct {
	Message Message
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
func (s *Sender) SendSMS(request ReqMessage) (ResponseBody, error) {
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
	// set channel type
	// 	0 : Normal SMS
	//  1 : Alert SMS
	//  2 : OTP SMS
	data.Set("channel", "2")

	// set proxy
	proxyUrl := s.Config.ProxyUrl
	if proxyUrl != "" {
		err := os.Setenv(CURLOPT_PROXY, proxyUrl)
		if err != nil {
			log.Error(err)
			return ResponseBody{}, err
		}
	}

	req, err := http.NewRequest("POST", urlPath, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error(err)
		return ResponseBody{}, err
	}

	req.Header.Set("Content-Type", "x-www-form-urlencoded")

	timeout := s.Config.ConnectTimeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	res, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return ResponseBody{}, err
	}

	if res.StatusCode >= 400 {
		log.Error(res)
		return ResponseBody{}, errors.New("failed to send SMS")
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	resBody := ResponseBody{}

	if err != nil {
		log.Error(res)
		return ResponseBody{}, err
	}

	status := strings.Replace(string(body), "Status=", "", 1)
	resBody.Message.To = request.PhoneNumber
	resBody.Message.Status = status

	return resBody, nil
}
