package sms_jatis

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/fairyhunter13/pool"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"net/url"
)

// Client is the interface of Jatis SMS client.
type Client interface {
	// SendSMS sends message to the Jatis platform.
	SendSMS(ctx context.Context, request *RequestMessage) (respBody *ResponseMessage, err error)
}

type client struct {
	opt *Option
}

func (c *client) Assign(o *Option) *client {
	if o == nil {
		return c
	}

	c.opt = o.Clone()
	return c
}

// NewClient initializes a new client with the given option.
func NewClient(opts ...FnOption) (c Client) {
	o := (new(Option)).Assign(opts...).Default()
	c = (new(client)).Assign(o)
	return
}

// List of form keys used for request in Jatis.
const (
	FormKeyUserID    = "userid"
	FormKeyPassword  = "password"
	FormKeySender    = "sender"
	FormKeyMSISDN    = "msisdn"
	FormKeyMessage   = "message"
	FormKeyDivision  = "division"
	FormKeyBatchName = "batchname"
	FormKeyUploadBy  = "uploadby"
	FormKeyChannel   = "channel"
)

func (c *client) getRequestFormData(req *RequestMessage) url.Values {
	data := url.Values{}
	data.Set(FormKeyUserID, c.opt.UserID)
	data.Set(FormKeyPassword, c.opt.Password)
	data.Set(FormKeySender, c.opt.Sender)
	data.Set(FormKeyDivision, c.opt.Division)
	data.Set(FormKeyUploadBy, c.opt.UploadBy)
	data.Set(FormKeyMSISDN, req.PhoneNumber)
	data.Set(FormKeyMessage, req.Text)
	data.Set(FormKeyBatchName, req.BatchName)
	data.Set(FormKeyChannel, req.getChannel())
	return data
}

func (c *client) getRequestDR(req *RequestDRPull) (io.Reader, error) {
	res, err := xml.Marshal(req.DRRequest)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(res), nil
}

func (c *client) setHeaders(req *http.Request) {
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationForm)
	for _, ipStr := range c.opt.CustomIPs {
		req.Header.Set(fiber.HeaderXForwardedFor, ipStr)
	}
}

func (c *client) writeRequest(buff *bytes.Buffer, request *RequestMessage) (err error) {
	data := c.getRequestFormData(request.Default(c.opt))
	_, err = buff.WriteString(data.Encode())
	return
}

func (c *client) prepareRequest(ctx context.Context, buff *bytes.Buffer) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, http.MethodPost, c.opt.BaseURL+c.opt.SendMessagePath, buff)
	if err != nil {
		return
	}

	c.setHeaders(req)
	return
}

func (c *client) prepareRequestXML(ctx context.Context, xmlReq io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(http.MethodPost, c.opt.BaseURL+c.opt.DeliveryReportPath, xmlReq)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	c.setHeaders(req)
	return
}

func (c *client) doRequest(ctx context.Context, request *RequestMessage) (resp *http.Response, err error) {
	buff := pool.GetBuffer()
	defer pool.Put(buff)

	err = c.writeRequest(buff, request)
	if err != nil {
		return
	}

	req, err := c.prepareRequest(ctx, buff)
	if err != nil {
		return
	}

	resp, err = c.opt.client.Do(req)
	return
}

func (c *client) doRequestDR(ctx context.Context, request *RequestDRPull) (resp *http.Response, err error) {
	xmlReq, err := c.getRequestDR(request)
	if err != nil {
		return
	}

	req, err := c.prepareRequestXML(ctx, xmlReq)
	if err != nil {
		return
	}

	resp, err = c.opt.client.Do(req)
	return
}

func (c *client) getResponseBody(resp *http.Response) (byteBody []byte, err error) {
	byteBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = getUnknownErr(resp, byteBody)
	return
}

// SendSMS sends message to the Jatis platform.
func (c *client) SendSMS(ctx context.Context, request *RequestMessage) (response *ResponseMessage, err error) {
	if request == nil {
		err = ErrNilArgs
		return
	}

	resp, err := c.doRequest(ctx, request)
	defer func() {
		if resp == nil || resp.Body == nil {
			return
		}

		_ = resp.Body.Close()
	}()
	if err != nil {
		return
	}

	body, err := c.getResponseBody(resp)
	if err != nil {
		return
	}

	respMap, err := url.ParseQuery(string(body))
	if err != nil {
		return
	}

	response, err = (new(ResponseMessage)).
		assign(respMap).
		getRespAndErr()
	return
}

// DeliveryReportPull get delivery report status send message
func (c *client) DeliveryReportPull(ctx context.Context, request *RequestDRPull) (response *ResponseDRPull, err error) {
	if request == nil {
		err = ErrNilArgs
		return
	}

	resp, err := c.doRequestDR(ctx, request)

	defer func() {
		if resp == nil || resp.Body == nil {
			return
		}

		_ = resp.Body.Close()
	}()

	if err != nil {
		return
	}

	err = xml.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return
	}

	return
}
