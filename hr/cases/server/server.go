package server

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/rm4n0s/collabeek/common"
	createmember "github.com/rm4n0s/collabeek/hr/cases/create_member"
	isauthenticated "github.com/rm4n0s/collabeek/hr/cases/is_authenticated"
	"github.com/rm4n0s/collabeek/hr/cases/login"
	"github.com/rm4n0s/collabeek/hr/cases/registration"
	"github.com/rm4n0s/collabeek/hr/db"
)

func NewHrCollabeekEchoServer(input ServerInput) (*echo.Echo, error) {
	ctx := context.Background()
	conn, err := sql.Open("postgres", input.DBHost)
	if err != nil {
		return nil, err
	}
	queries := db.New(conn)
	defer conn.Close()

	smtpService := common.NewSmtpService(input.SenderEmail,
		input.SmtpUsername,
		input.SmtpPassword,
		input.SmtpHost,
		input.SmtpPort)

	initializeAdmin(ctx, queries, smtpService, input.AdminEmail)

	e := echo.New()
	e.Use(echojwt.JWT([]byte(input.TokenSecret)))
	e.Validator = common.NewValidator()
	lh := login.NewLoginHandler(queries, input.TokenSecret)
	ah := isauthenticated.NewIsAuthenticatedHandler(input.TokenSecret)
	cmh := createmember.NewCreateMemberHandler(queries, smtpService)
	rh := registration.NewRegistrationHandler(queries, smtpService)
	g := e.Group("/api/v1")
	g.POST("/login", lh.LoginHandler)
	g.POST("/is_authenticated", ah.IsAuthenticatedHandler)
	g.POST("/members", cmh.CreateMemberHandler, common.OnlyAdminOrModerators)
	g.POST("/registration", rh.RegistrationHandler)
	return e, nil
}
