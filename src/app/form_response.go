package app

import "github.com/linkernetworks/oauth/src/validator"

type FormActionResponse struct {
	Error       bool                    `json:"error"`
	Message     string                  `json:"message"`
	Validations validator.ValidationMap `json:"validations"`
}
