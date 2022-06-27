package sms_jatis

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

const (
	baseURL = "https://sms-api.jatismobile.com/index.ashx"
)

var (
	cfg = Config{
		BaseUrl: baseURL,
	}
)

func TestSender_SendSMS(t *testing.T) {
	type args struct {
		ctx     context.Context
		request ReqMessage
	}
	type jatisResponse struct {
		code int
		body string
		err  error
	}
	ctx := context.Background()
	tests := []struct {
		name           string
		wantRespBody   ResponseBody
		wantErr        bool
		config         Config
		args           args
		clientResponse jatisResponse
	}{
		{
			name: "Given successfull sending message",
			args: args{
				request: ReqMessage{
					PhoneNumber: "+0912122",
					Text:        "123412",
				},
				ctx: ctx,
			},
			config: cfg,
			clientResponse: jatisResponse{
				code: http.StatusOK,
				body: "Status=1&MessageId=8989",
			},
			wantRespBody: ResponseBody{
				MessageId: "8989",
				To:        "+0912122",
				Status:    "1",
			},
		},
		{
			name:    "Given error httprequest context",
			wantErr: true,
			args: args{
				request: ReqMessage{
					PhoneNumber: "+0912122",
					Text:        "123412",
				},
				ctx: nil,
			},
			config:         cfg,
			clientResponse: jatisResponse{},
		},
		{
			name:    "Given error client do",
			wantErr: true,
			args: args{
				request: ReqMessage{
					PhoneNumber: "+0911",
					Text:        "123412",
				},
				ctx: ctx,
			},
			config: Config{},
			clientResponse: jatisResponse{
				code: http.StatusOK,
				body: "Status=2",
			},
		},
		{
			name:    "Given response not found",
			wantErr: true,
			args: args{
				request: ReqMessage{
					PhoneNumber: "+0912122",
					Text:        "123412",
				},
				ctx: ctx,
			},
			config: cfg,
			clientResponse: jatisResponse{
				code: http.StatusNotFound,
			},
		},
		{
			name:    "Given error status",
			wantErr: true,
			args: args{
				request: ReqMessage{
					PhoneNumber: "+0912122",
					Text:        "123412",
				},
				ctx: ctx,
			},
			config: cfg,
			clientResponse: jatisResponse{
				code: http.StatusOK,
				body: "Status=2",
			},
			wantRespBody: ResponseBody{
				To:     "+0912122",
				Status: "missing Parameter",
			},
		},
		{
			name:    "Given error parse query string",
			wantErr: true,
			args: args{
				request: ReqMessage{
					PhoneNumber: "+0912122",
					Text:        "123412",
				},
				ctx: ctx,
			},
			config: cfg,
			clientResponse: jatisResponse{
				code: http.StatusOK,
				body: "//%%mailto",
			},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.config)

			httpmock.RegisterResponder(http.MethodPost, cfg.BaseUrl,
				func(req *http.Request) (*http.Response, error) {
					return httpmock.NewStringResponse(tt.clientResponse.code, tt.clientResponse.body), tt.clientResponse.err
				})

			gotRespBody, err := s.SendSMS(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sender.SendSMS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRespBody, tt.wantRespBody) {
				t.Errorf("Sender.SendSMS() = %v, want %v", gotRespBody, tt.wantRespBody)
			}
		})
	}
}
