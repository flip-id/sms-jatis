package sms_jatis

var JatisError = map[int]string{
	ErrorMissingParameter:    "missing Parameter",
	ErrorInvalidUserPassword: "invalid user or password",
	ErrorInvalidMessage:      "invalid message",
	ErrorRecepientNumber:     "invalid recipient number",
	ErrorInvalidSender:       "invalid sender name",
	ErrorIPNotAllowed:        "ip address not allowed",
	ErrorInternal:            "internal server error",
	ErrorInvalidDivision:     "invalid division name",
	ErrorInvalidChannel:      "invalid channel type",
	ErrorInsufficientToken:   "insufficient token",
	ErrorTokenNotAvailable:   "token not available",
	ErrorGeneral:             "unknown error",
}

const (
	ErrorGeneral             = 0
	StatusSuccess            = 1
	ErrorMissingParameter    = 2
	ErrorInvalidUserPassword = 3
	ErrorInvalidMessage      = 4
	ErrorRecepientNumber     = 5
	ErrorInvalidSender       = 6
	ErrorIPNotAllowed        = 7
	ErrorInternal            = 8
	ErrorInvalidDivision     = 9
	ErrorInvalidChannel      = 20
	ErrorInsufficientToken   = 21
	ErrorTokenNotAvailable   = 22
)
