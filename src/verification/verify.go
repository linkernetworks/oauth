package verification

import (
	"github.com/linkernetworks/oauth/src/sms"
	"errors"
	"fmt"
)

var SmsVerificationCodeMessageFormat string = " Your verification code is %s"

type VerificationProcessReceiver interface {
	SetVerificationCode(code string)

	GetVerificationCode() string

	SetVerified(val bool) bool

	MatchVerificationCode(code string) bool

	GetCellPhoneNumber() string
}

func Verify(target VerificationProcessReceiver, code string) bool {
	return target.SetVerified(target.GetVerificationCode() == code)
}

func Send(sms *sms.SMSClient, target VerificationProcessReceiver, codeLength int) error {
	to := target.GetCellPhoneNumber()
	if to == "" {
		return errors.New("cellphone number is empty")
	}

	code := GenerateCode(codeLength)
	target.SetVerificationCode(code)
	msg := fmt.Sprintf(SmsVerificationCodeMessageFormat, code)
	_, _, err := sms.SendMessage(to, msg)
	return err
}
