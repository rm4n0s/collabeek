package server

import (
	"context"

	"github.com/rm4n0s/collabeek/hr/db"
)

type ServerInput struct {
	AdminEmail   string
	SenderEmail  string
	SmtpUsername string
	SmtpPassword string
	SmtpHost     string
	SmtpPort     int
	DBHost       string
	TokenSecret  string
}

type InitializationDB interface {
	CreateMember(ctx context.Context, arg db.CreateMemberParams) (db.Member, error)
	GetMemberByEmail(ctx context.Context, email string) (db.Member, error)
}
