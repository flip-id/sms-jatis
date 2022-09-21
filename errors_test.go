package sms_jatis

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUnknownError_Error(t *testing.T) {
	type fields struct {
		StatusCode int
		Message    []byte
	}
	tests := []struct {
		name   string
		fields func() fields
		want   func() string
	}{
		{
			name: "unknown error - status unauthorized",
			fields: func() fields {
				return fields{
					StatusCode: http.StatusUnauthorized,
					Message:    []byte(http.StatusText(http.StatusUnauthorized)),
				}
			},
			want: func() string {
				return fmt.Sprintf(
					"unknown error SMS Jatis: status code:%d message:%s",
					http.StatusUnauthorized,
					[]byte(http.StatusText(http.StatusUnauthorized)),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			e := UnknownError{
				StatusCode: fields.StatusCode,
				Message:    fields.Message,
			}
			assert.Equalf(t, tt.want(), e.Error(), "Error()")
		})
	}
}
