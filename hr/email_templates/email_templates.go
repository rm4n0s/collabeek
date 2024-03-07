package emailtemplates

import _ "embed"

type RegistrationLetterInput struct {
	RegistrationCode string
}

type SuccessfulRegistrationLetterInput struct {
	Fullname string
}

//go:embed registration_letter.html
var RegistrationLetter string

//go:embed successful_registration_letter.html
var SuccessfulRegistrationLetter string
