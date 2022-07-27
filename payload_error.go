package sms_jatis

import "fmt"

// List of all status parameters used in Jatis.
const (
	StatusSuccess uint64 = iota + 1
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
	StatusInvalidChannel uint64 = iota + 20
	StatusTokenNotEnough
	StatusTokenNotAvailable
)

var statusMapping = map[uint64]string{
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
