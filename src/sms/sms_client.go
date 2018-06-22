package sms

import "github.com/sfreiberg/gotwilio"

// TwilioConfig
type TwilioConfig struct {
	Sid       string `json:"sid"`
	AuthToken string `json:"token"`
	CallFrom  string `json:"callfrom"`
}

// SMSClient :
type SMSClient struct {
	ID    string
	Token string

	Sid       string
	AuthToken string
	CallFrom  string
}

func NewSMSClient(id, pw, from string) *SMSClient {
	return &SMSClient{
		Sid:       id,
		AuthToken: pw,
		CallFrom:  from}
}

func NewSMSClientFromConfig(twilioConfig TwilioConfig) *SMSClient {
	return &SMSClient{
		Sid:       twilioConfig.Sid,
		AuthToken: twilioConfig.AuthToken,
		CallFrom:  twilioConfig.CallFrom,
	}
}

// SMS is send sms to specific number.
func (s *SMSClient) SendMessage(to, msg string) (*gotwilio.SmsResponse, *gotwilio.Exception, error) {
	client := gotwilio.NewTwilioClient(s.Sid, s.AuthToken)
	return client.SendSMS(s.CallFrom, to, msg, "", "")
}
