package sms_jatis

import (
	"context"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient(
		WithBaseURL(DefaultBaseURL),
		WithDivision("Flip-Communication"),
		WithUserID("flip"),
		WithPassword("flip"),
		WithSender("Flip"),
		WithUploadBy("Flip"),
		WithChannel(ChannelOTP),
		WithCustomIPs("34.134.35.34"),
		WithTimeout(DefaultTimeout),
		WithClient(http.DefaultClient),
		WithHystrixOptions(),
	)
	assert.NotNil(t, c)
}

func Test_client_SendSMS(t *testing.T) {
	type fields struct {
		opt *Option
	}
	type args struct {
		ctx     context.Context
		request *RequestMessage
	}
	tests := []struct {
		name         string
		fields       func() fields
		args         func() args
		wantResponse func() *ResponseMessage
		wantErr      bool
	}{
		{
			name: "nil args",
			fields: func() fields {
				return fields{}
			},
			args: func() args {
				return args{
					ctx:     context.TODO(),
					request: nil,
				}
			},
			wantResponse: func() *ResponseMessage {
				return nil
			},
			wantErr: true,
		},
		{
			name: "success request",
			fields: func() fields {
				client := new(http.Client)
				httpmock.ActivateNonDefault(client)
				opt := (new(Option)).Assign(
					WithBaseURL(DefaultBaseURL),
					WithDivision("Flip-Communication"),
					WithUserID("flip"),
					WithPassword("flip"),
					WithSender("Flip"),
					WithUploadBy("Flip"),
					WithCustomIPs("34.134.35.34"),
					WithTimeout(DefaultTimeout),
					WithClient(client),
					WithChannel(ChannelOTP),
				).Default()
				respByte := []byte("Status=1&MessageId=3120910074119080761033f9")
				httpmock.RegisterResponder(
					http.MethodPost,
					DefaultBaseURL,
					httpmock.NewBytesResponder(http.StatusOK, respByte),
				)
				return fields{
					opt: opt,
				}
			},
			args: func() args {
				chanOTP := ChannelOTP
				req := &RequestMessage{
					PhoneNumber: "628123123123",
					Text:        "Hello!",
					BatchName:   "Flip",
					Channel:     &chanOTP,
				}
				return args{
					ctx:     context.TODO(),
					request: req,
				}
			},
			wantResponse: func() *ResponseMessage {
				return &ResponseMessage{
					MessageID:    "3120910074119080761033f9",
					Status:       StatusSuccess.String(),
					StatusNumber: StatusSuccess,
				}
			},
			wantErr: false,
		},
		{
			name: "missing user ID",
			fields: func() fields {
				client := new(http.Client)
				httpmock.ActivateNonDefault(client)
				opt := (new(Option)).Assign(
					WithBaseURL(DefaultBaseURL),
					WithDivision("Flip-Communication"),
					WithPassword("flip"),
					WithSender("Flip"),
					WithUploadBy("Flip"),
					WithCustomIPs("34.134.35.34"),
					WithTimeout(DefaultTimeout),
					WithClient(client),
					WithChannel(ChannelOTP),
				).Default()
				respByte := []byte("Status=2")
				httpmock.RegisterResponder(
					http.MethodPost,
					DefaultBaseURL,
					httpmock.NewBytesResponder(http.StatusOK, respByte),
				)
				return fields{
					opt: opt,
				}
			},
			args: func() args {
				chanOTP := ChannelOTP
				req := &RequestMessage{
					PhoneNumber: "628123123123",
					Text:        "Hello!",
					BatchName:   "Flip",
					Channel:     &chanOTP,
				}
				return args{
					ctx:     context.TODO(),
					request: req,
				}
			},
			wantResponse: func() *ResponseMessage {
				return nil
			},
			wantErr: true,
		},
		{
			name: "unknown error",
			fields: func() fields {
				client := new(http.Client)
				httpmock.ActivateNonDefault(client)
				opt := (new(Option)).Assign(
					WithClient(client),
					WithChannel(ChannelOTP),
					WithCustomIPs("34.35.134.34"),
				).Default()
				respByte := []byte(http.StatusText(http.StatusUnauthorized))
				httpmock.RegisterResponder(
					http.MethodPost,
					DefaultBaseURL,
					httpmock.NewBytesResponder(http.StatusUnauthorized, respByte),
				)
				return fields{
					opt: opt,
				}
			},
			args: func() args {
				chanOTP := ChannelOTP
				req := &RequestMessage{
					PhoneNumber: "628123123123",
					Text:        "Hello!",
					BatchName:   "Flip",
					Channel:     &chanOTP,
				}
				return args{
					ctx:     context.TODO(),
					request: req,
				}
			},
			wantResponse: func() *ResponseMessage {
				return nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer httpmock.DeactivateAndReset()
			c := &client{
				opt: tt.fields().opt,
			}
			args := tt.args()
			gotResponse, err := c.SendSMS(args.ctx, args.request)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equalf(t, tt.wantResponse(), gotResponse, "SendSMS(%v, %v)", args.ctx, args.request)
		})
	}
}
