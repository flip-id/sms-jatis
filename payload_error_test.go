package sms_jatis

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatusParam_Uint64String(t *testing.T) {
	tests := []struct {
		name string
		s    StatusParam
		want string
	}{
		{
			name: "success",
			s:    StatusSuccess,
			want: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.s.Uint64String(), "Uint64String()")
		})
	}
}

func TestResponseMessage_Error(t *testing.T) {
	type fields struct {
		MessageID    string
		Status       string
		StatusNumber StatusParam
	}
	tests := []struct {
		name   string
		fields func() fields
		want   func() string
	}{
		{
			name: "missing parameter",
			fields: func() fields {
				return fields{
					MessageID:    "3120910074119080761033f9",
					Status:       StatusMissingParameter.String(),
					StatusNumber: StatusMissingParameter,
				}
			},
			want: func() string {
				return fmt.Sprintf(
					"error response from SMS Jatis, status:%d status text:%s",
					StatusMissingParameter,
					StatusMissingParameter.String(),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			r := &ResponseMessage{
				MessageID:    fields.MessageID,
				Status:       fields.Status,
				StatusNumber: fields.StatusNumber,
			}
			assert.Equalf(t, tt.want(), r.Error(), "Error()")
		})
	}
}
