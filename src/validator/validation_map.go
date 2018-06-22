package validator

type ValidationMap map[string]FieldValidation

func (v ValidationMap) HasError() bool {
	for _, it := range v {
		if it.Error {
			return true
		}
	}
	return false
}

type FieldValidation struct {
	Field   string `json:"field"`
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
