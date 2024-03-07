package login

import (
	"context"
	"time"

	"github.com/rm4n0s/collabeek/hr/db"
)

type LoginDB interface {
	GetMemberByUsernameAndPassword(ctx context.Context, arg db.GetMemberByUsernameAndPasswordParams) (db.Member, error)
}

type LoginHandler struct {
	db          LoginDB
	tokenSecret []byte
}

type LoginForm struct {
	Username    string        `validate:"required"`
	Password    string        `validate:"required"`
	ExpireAfter time.Duration `validate:"required,gt=1"`
}

type LoginResponse struct {
	ID       int32
	Username string
	Role     string
	Token    string
}
