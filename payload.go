package sms_jatis

import (
	"github.com/fairyhunter13/phone"
	"github.com/fairyhunter13/reflecthelper/v5"
	"net/url"
	"strconv"
	"time"
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

// RequestDRPull defines the request Delivery Report to Jatis.
type RequestDRPull struct {
	DRRequest RequestDR `xml:"DRRequest"`
}

type RequestDR struct {
	UserID    string   `xml:"UserId"`
	Password  string   `xml:"Password"`
	Sender    string   `xml:"Sender"`
	MessageID []string `xml:"MessageId"`
}

// ResponseDRPull defines the Delivery Report Pull from Jatis.
type ResponseDRPull struct {
	DRResponse
}

// DRResponse
// RequestStatus codes list:
// 	6 = Missing Parameter
// 	7 = Invalid User Id or Password
// 	8 = Invalid Sender
// 	9 = Client’s IP Address is not allowed
//	10 = Internal Server Error
// DeliveryStatus codes list:
// 	1 = Success, received -> Success/message processed successfully by Telco.
// 	2 = Success, not received -> The SMS sent successfully to Telco but not received
// 	3 = Success, unknown number -> The Destination Number is not valid
// 	4 = Failed -> SMS submission failed
// 	77 = Delivery Status is not available -> The Delivery Status of requested Message Id currently not available
//			Delivery Status available in 1 to 3 days after SMS submission.
type DRResponse struct {
	ClientID           string
	Sender             string
	IsEncrypted        bool
	RequestStatus      int
	JatisMessageID     string
	DeliveryStatus     int
	DateStatusReceived time.Time
}

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
