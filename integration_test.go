//go:build integration
// +build integration

package sms_jatis

import (
	"context"
	"flag"
	"fmt"
	"github.com/fairyhunter13/dotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	// DefaultBitSizeUint is the default bit size of uint.
	DefaultBitSizeUint = 64
)

var (
	c    Client
	once sync.Once
)

func setupClient() {
	once.Do(func() {
		err := dotenv.Load2(
			dotenv.WithPaths(".env"),
		)
		if err != nil {
			log.Fatalln(err)
		}

		channelSMS, err := strconv.ParseUint(os.Getenv("SMS_JATIS_CHANNEL"), DefaultBaseDecimal, DefaultBitSizeUint)
		if err != nil {
			log.Fatalln(err)
		}

		c = NewClient(
			WithUserID(os.Getenv("SMS_JATIS_USER_ID")),
			WithPassword(os.Getenv("SMS_JATIS_PASSWORD")),
			WithSender(os.Getenv("SMS_JATIS_SENDER")),
			WithDivision(os.Getenv("SMS_JATIS_DIVISION")),
			WithUploadBy(os.Getenv("SMS_JATIS_UPLOAD_BY")),
			WithChannel(channelSMS),
			WithCustomIPs(os.Getenv("SMS_JATIS_CUSTOM_IP")),
		)
		if err != nil {
			log.Fatalln(err)
		}
	})
}

// Run integration tests.
// Notes: Run this test only on local, not on CI/CD.
func TestMain(m *testing.M) {
	flag.Parse()
	setupClient()

	os.Exit(m.Run())
}

func formatTime(in time.Time) string {
	return in.Format("2006-01-02 15:04:05")
}

func TestSendMessageSuccessful(t *testing.T) {
	ctx := context.Background()
	req := &RequestMessage{
		PhoneNumber: os.Getenv("PHONE_NUMBER"),
		Text:        fmt.Sprintf("Hello, the timestamp is %s. [Jatis Integration Testing]", formatTime(time.Now())),
	}
	resp, err := c.SendSMS(ctx, req)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, StatusSuccess, resp.StatusNumber)
	assert.Equal(t, statusMapping[StatusSuccess], resp.Status)
	assert.NotEmpty(t, resp.MessageID)
}

func TestSendMessageFailed(t *testing.T) {
	t.Parallel()
	t.Run("no phone number", func(t *testing.T) {
		ctx := context.Background()
		req := &RequestMessage{
			PhoneNumber: "",
			Text:        fmt.Sprintf("Hello, the timestamp is %s. [Jatis Integration Testing]", formatTime(time.Now())),
		}
		resp, err := c.SendSMS(ctx, req)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.Equal(t, reflect.TypeOf(new(ResponseMessage)), reflect.TypeOf(err))
		resp = err.(*ResponseMessage)
		assert.Equal(t, StatusMissingParameter, resp.StatusNumber)
		assert.Equal(t, statusMapping[StatusMissingParameter], resp.Status)
		assert.Empty(t, resp.MessageID)
	})

}
