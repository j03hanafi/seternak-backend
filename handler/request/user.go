package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// SignUp defines the request payload for SignUp method.
type SignUp struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// Validate validates the SignUp request fields.
func (s SignUp) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Email, validation.Required, is.Email),
		validation.Field(&s.Password, validation.Required, validation.Length(8, 255)),
		validation.Field(&s.Name, validation.Required, validation.Length(1, 255)),
	)
}
