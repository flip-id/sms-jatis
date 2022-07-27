package sms_jatis

import (
	"github.com/fairyhunter13/phone"
	"github.com/fairyhunter13/reflecthelper/v5"
	"net/url"
	"strconv"
)

// List of channel used in Jatis.
const (
	ChannelRegular uint64 = iota
	ChannelAlert
	ChannelOTP
)

const (
	// DefaultBatchOTP is the default batch name for OTP.
	DefaultBatchOTP = "otp"
	// DefaultBaseDecimal is the default base 10 for decimal.
	DefaultBaseDecimal = 10
)

// ResponseMessage defines the response message from Jatis.
type ResponseMessage struct {
	MessageID    string
	Status       string
	StatusNumber uint64
}

// List of key used to access the value of URL decoded values.
const (
	KeyURLDecodedStatus    = "Status"
	KeyURLDecodedMessageID = "MessageId"
)

func (r *ResponseMessage) assign(val url.Values) *ResponseMessage {
	r.MessageID = val.Get(KeyURLDecodedMessageID)
	r.StatusNumber = reflecthelper.GetUint(val.Get(KeyURLDecodedStatus))
	r.Status = statusMapping[r.StatusNumber]
	return r
}

// RequestMessage defines the request message to Jatis.
type RequestMessage struct {
	PhoneNumber string
	Text        string
	BatchName   string
	Channel     *uint64
}

func (r *RequestMessage) getChannel() string {
	if r.Channel == nil {
		return ""
	}

	return strconv.FormatUint(*r.Channel, DefaultBaseDecimal)
}

func getChannelOTP() *uint64 {
	chanOTP := ChannelOTP
	return &chanOTP
}

// Default returns the default request message for Jatis.
func (r *RequestMessage) Default(o *Option) *RequestMessage {
	if r.BatchName == "" {
		r.BatchName = DefaultBatchOTP
	}

	if r.Channel == nil || *r.Channel > ChannelOTP {
		if o == nil || o.Channel == nil {
			r.Channel = getChannelOTP()
			goto next
		}

		r.Channel = o.Channel
	}

next:

	r.PhoneNumber = phone.NormalizeID(r.PhoneNumber, 0)
	return r
}
