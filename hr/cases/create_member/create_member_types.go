package createmember

import (
	"context"
	"time"

	"github.com/rm4n0s/collabeek/common"
	"github.com/rm4n0s/collabeek/hr/db"
)

type CreateMemberDB interface {
	GetMemberByEmail(ctx context.Context, email string) (db.Member, error)
	CreateMember(ctx context.Context, arg db.CreateMemberParams) (db.Member, error)
}

type CreateMemberHandler struct {
	smtpService common.ISmtpService
	db          CreateMemberDB
}

type CreateMemberForm struct {
	Email string `validate:"required,email"`
	Role  string `validate:"required,oneof=admin moderator member"`
}

type CreateMemberResponse struct {
	ID        int32
	Email     string
	Role      string
	CreatedAt time.Time
}
