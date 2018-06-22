package sms

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

const TEST_MESSAGE = "test message from oauth/sms package"

func newTwilioConfig() TwilioConfig {
	configPath := os.Getenv("OAUTH_CONFIG_PATH")
	if configPath == "" {
		return TwilioConfig{}
	}

	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Fatalf("read config file %s error: %v\n", configPath, err)
	}

	var twilioConfig TwilioConfig
	err = json.Unmarshal(content, &twilioConfig)
	if err != nil {
		logrus.Fatalf("parse twilio config error: %v\n", err)
	}

	return twilioConfig
}

func TestSendSMSMessage(t *testing.T) {
	twilioConfig := newTwilioConfig()
	if twilioConfig.Sid == "" {
		t.Skip()
	}

	phoneNumber, found := os.LookupEnv("TEST_SMS_PHONENUMBRR")
	if !found {
		t.Logf("TEST_SMS_PHONENUMBRR is undefined")
		return
	}

	t.Logf("Sending SMS to phone number: %s", phoneNumber)
	sms := NewSMSClientFromConfig(twilioConfig)
	resp, exception, err := sms.SendMessage(phoneNumber, TEST_MESSAGE)

	if exception != nil {
		t.Errorf("send sms error: %v", exception)
	}

	if err == nil {
		t.Logf("Successfully sending SMS, response: %v", resp)
	} else {
		t.Errorf("Sending SMS failed. Please see https://www.twilio.com/console/sms/logs for more details.\n Error:%v\n Exception: %v", err, exception)
	}
}
