package app

import "bitbucket.org/linkernetworks/cv-tracker/src/oauth/validator"

type FormActionResponse struct {
	Error       bool                    `json:"error"`
	Message     string                  `json:"message"`
	Validations validator.ValidationMap `json:"validations"`
}
