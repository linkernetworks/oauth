package verification

import "testing"

func TestGenerateVerifyCode(t *testing.T) {
	code := GenerateCode(0)
	if code != "" {
		t.Errorf("verification code generation failed")
	}

}

func TestGenerateVerifyCodeFor6Digit(t *testing.T) {
	code := GenerateCode(6)
	if code == "" || len(code) != 6 {
		t.Errorf("verification code generation failed")
	}
}
