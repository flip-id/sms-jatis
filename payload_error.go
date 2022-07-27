package sms_jatis

import "fmt"

// StatusParam is the status parameter in the response message.
type StatusParam uint64

// String returns the string representation of the status param.
func (s StatusParam) String() string {
	return statusMapping[s]
}

// List of all status parameters used in Jatis.
const (
	StatusSuccess StatusParam = iota + 1
	StatusMissingParameter
	StatusInvalidUserIDOrPassword
	StatusInvalidMessage
	StatusInvalidMSISDN
	StatusInvalidSender
	StatusDeniedIPAddress
	StatusInternalServerError
	StatusInvalidDivision
)

// List of all status parameters used in Jatis.
const (
	StatusInvalidChannel StatusParam = iota + 20
	StatusTokenNotEnough
	StatusTokenNotAvailable
)

var statusMapping = map[StatusParam]string{
	StatusSuccess:                 "Success",
	StatusMissingParameter:        "Missing Parameter",
	StatusInvalidUserIDOrPassword: "Invalid User Id or Password",
	StatusInvalidMessage:          "Invalid Message",
	StatusInvalidMSISDN:           "Invalid MSISDN",
	StatusInvalidSender:           "Invalid Sender",
	StatusDeniedIPAddress:         "Client’s IP Address is not allowed",
	StatusInternalServerError:     "Internal Server Error",
	StatusInvalidDivision:         "Invalid Division",
	StatusInvalidChannel:          "Invalid Channel",
	StatusTokenNotEnough:          "Token Not Enough",
	StatusTokenNotAvailable:       "Token Not Available",
}

func (r *ResponseMessage) getRespAndErr() (*ResponseMessage, error) {
	if r.StatusNumber == StatusSuccess {
		return r, nil
	}

	return nil, r
}

// ResponseMessage implements the error interface.
func (r *ResponseMessage) Error() string {
	return fmt.Sprintf(
		"error response from SMS Jatis, status:%d status text:%s",
		r.StatusNumber,
		r.Status,
	)
}
