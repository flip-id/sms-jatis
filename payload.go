package sms_jatis

import (
	"github.com/fairyhunter13/phone"
	"github.com/fairyhunter13/reflecthelper/v5"
	"net/url"
	"strconv"
)

// ChannelType specifies the channel type used in this package.
type ChannelType uint64

// List of channel used in Jatis.
const (
	ChannelRegular ChannelType = iota
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
	StatusNumber StatusParam
}

// List of key used to access the value of URL decoded values.
const (
	KeyURLDecodedStatus    = "Status"
	KeyURLDecodedMessageID = "MessageId"
)

func (r *ResponseMessage) assign(val url.Values) *ResponseMessage {
	r.MessageID = val.Get(KeyURLDecodedMessageID)
	r.StatusNumber = StatusParam(reflecthelper.GetUint(val.Get(KeyURLDecodedStatus)))
	r.Status = r.StatusNumber.String()
	return r
}

// RequestMessage defines the request message to Jatis.
type RequestMessage struct {
	PhoneNumber string
	Text        string
	BatchName   string
	Channel     *ChannelType
}

func (r *RequestMessage) getChannel() string {
	if r.Channel == nil {
		return ""
	}

	return strconv.FormatUint(uint64(*r.Channel), DefaultBaseDecimal)
}

func getChannelOTP() *ChannelType {
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
