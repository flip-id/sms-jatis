package sms_jatis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

// ReqJatis request for Jatis
type ReqJatis struct {
	UserId    string `json:"userid"`
	Password  string `json:"password"`
	Msisdn    string `json:"msisdn"`
	Message   string `json:"message"`
	Sender    string `json:"sender"`
	Batchname string `json:"batchname"`
	Division  string `json:"division"`
	UploadBy  string `json:"uploadby"`
	Channel   int    `json:"channel"`
}

// end of request

func New(config Config) *Sender {
	return &Sender{
		Config: config,
	}
}

// CallbackData callback data
type CallbackData struct {
	Error struct {
		Description string `json:"description"`
		GroupID     int    `json:"group_id"`
		GroupName   string `json:"group_name"`
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Permanent   bool   `json:"permanent"`
	} `json:"error"`
	To     string `json:"to"`
	Status int    `json:"status"`
	SentAt string `json:"sent_at"`
}

// end of callback data

// SendSMS function to send message
func (s *Sender) SendSMS(request ReqMessage) (ResponseBody, error) {
	url := fmt.Sprintf("%s", s.Config.BaseUrl)
	payload, err := json.Marshal(ReqJatis{
		UserId:    s.Config.UserId,
		Password:  s.Config.Password,
		Msisdn:    request.PhoneNumber,
		Message:   request.Text,
		Sender:    s.Config.Sender,
		Batchname: "otp",
		Division:  s.Config.Division,
		UploadBy:  s.Config.UploadBy,
		Channel:   s.Config.channel,
	})

	if err != nil {
		log.Error(err)
		return ResponseBody{}, err
	}

	// set proxy
	proxyUrl := s.Config.ProxyUrl
	if proxyUrl != "" {
		err = os.Setenv(CURLOPT_PROXY, proxyUrl)
		if err != nil {
			log.Error(err)
			return ResponseBody{}, err
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))

	if err != nil {
		log.Error(err)
		return ResponseBody{}, err
	}

	req.Header.Set("Content-Type", "application/json")

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
