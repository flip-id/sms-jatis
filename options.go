package sms_jatis

import (
	"github.com/fairyhunter13/reflecthelper/v5"
	"github.com/gojek/heimdall/v7"
	"github.com/gojek/heimdall/v7/hystrix"
	"net/http"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the default base URL for the Jatis API.
	DefaultBaseURL = "https://sms-api.jatismobile.com/index.ashx"
	// DefaultTimeout is the default timeout for the Jatis API.
	DefaultTimeout = 30 * time.Second
)

// Option is a struct containing all configurations for Jatis.
type Option struct {
	BaseURL            string
	DeliveryReportPath string
	SendMessagePath    string
	UserID             string
	Password           string
	Sender             string
	Division           string
	UploadBy           string
	Channel            *ChannelType
	CustomIPs          []string
	Timeout            time.Duration
	Client             heimdall.Doer
	HystrixOptions     []hystrix.Option
	client             *hystrix.Client
}

// FnOption is a function that modifies Option.
type FnOption func(o *Option)

// WithBaseURL sets the base URL for the Jatis API.
func WithBaseURL(baseURL string) FnOption {
	return func(o *Option) {
		o.BaseURL = baseURL
	}
}

// WithUserID sets the user ID for the Jatis API.
func WithUserID(userID string) FnOption {
	return func(o *Option) {
		o.UserID = userID
	}
}

// WithPassword sets the password for the Jatis API.
func WithPassword(password string) FnOption {
	return func(o *Option) {
		o.Password = password
	}
}

// WithSender sets the sender for the Jatis API.
func WithSender(sender string) FnOption {
	return func(o *Option) {
		o.Sender = sender
	}
}

// WithDivision sets the division for the Jatis API.
func WithDivision(division string) FnOption {
	return func(o *Option) {
		o.Division = division
	}
}

// WithUploadBy sets the upload by for the Jatis API.
func WithUploadBy(uploadBy string) FnOption {
	return func(o *Option) {
		o.UploadBy = uploadBy
	}
}

// WithChannel sets the channel for the Jatis API.
func WithChannel(channel ChannelType) FnOption {
	return func(o *Option) {
		o.Channel = &channel
	}
}

// WithCustomIPs sets the custom IPs for the Jatis API.
func WithCustomIPs(customIPs ...string) FnOption {
	return func(o *Option) {
		o.CustomIPs = customIPs
	}
}

// WithTimeout sets the timeout for the Jatis API.
func WithTimeout(timeout time.Duration) FnOption {
	return func(o *Option) {
		o.Timeout = timeout
	}
}

// WithClient sets the client for the Jatis API.
func WithClient(client heimdall.Doer) FnOption {
	return func(o *Option) {
		o.Client = client
	}
}

// WithHystrixOptions sets the hystrix options for the Jatis API.
func WithHystrixOptions(options ...hystrix.Option) FnOption {
	return func(o *Option) {
		o.HystrixOptions = options
	}
}

// Clone returns a shallow clone of the Option.
func (o *Option) Clone() *Option {
	opt := *o
	return &opt
}

// Assign assigns the Option to the given FnOption.
func (o *Option) Assign(opts ...FnOption) *Option {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Default returns a default Option.
func (o *Option) Default() *Option {
	if o.BaseURL == "" {
		o.BaseURL = DefaultBaseURL
	}

	o.BaseURL = strings.TrimRight(o.BaseURL, "/")
	if o.Timeout < DefaultTimeout {
		o.Timeout = DefaultTimeout
	}

	if o.Channel == nil || *o.Channel > ChannelOTP {
		o.Channel = getChannelOTP()
	}

	if reflecthelper.IsNil(o.Client) {
		o.Client = http.DefaultClient
	}

	o.client = hystrix.NewClient(
		append([]hystrix.Option{
			hystrix.WithHystrixTimeout(o.Timeout),
			hystrix.WithHTTPTimeout(o.Timeout),
			hystrix.WithHTTPClient(o.Client),
		}, o.HystrixOptions...)...,
	)
	return o
}
