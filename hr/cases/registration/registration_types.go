package registration

import (
	"context"

	"github.com/rm4n0s/collabeek/common"
	"github.com/rm4n0s/collabeek/hr/db"
)

type RegistrationDB interface {
	GetMemberByEmail(ctx context.Context, email string) (db.Member, error)
	UpdateMemberForRegistration(ctx context.Context, arg db.UpdateMemberForRegistrationParams) error
	GetMemberByUsername(ctx context.Context, username string) (db.Member, error)
}

type RegistrationHandler struct {
	smtpService common.ISmtpService
	db          RegistrationDB
}

type RegistrationForm struct {
	Email              string `validate:"required,email"`
	Username           string `validate:"required,alphanum,min=4"`
	Fullname           string `validate:"required,min=4"`
	RegistrationSecret string `validate:"required,len=60"`
	Password           string `validate:"required,min=8"`
}
