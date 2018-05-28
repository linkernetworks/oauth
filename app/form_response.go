package app

import "github.com/linkernetworks/oauth/validator"

type FormActionResponse struct {
	Error       bool                    `json:"error"`
	Message     string                  `json:"message"`
	Validations validator.ValidationMap `json:"validations"`
}
