package sms_jatis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestMessage_Default(t *testing.T) {
	r := (new(RequestMessage)).Default(nil)
	assert.NotEmpty(t, r.BatchName)
	assert.NotEmpty(t, r.Channel)
}
