package login

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rm4n0s/collabeek/servers/common"
	"github.com/rm4n0s/collabeek/servers/hr/db"
)

func NewLoginHandler(db LoginDB, tokenSecret string) *LoginHandler {
	return &LoginHandler{
		db:          db,
		tokenSecret: []byte(tokenSecret),
	}
}

func (h *LoginHandler) LoginHandler(c echo.Context) error {
	form := new(LoginForm)
	if err := c.Bind(form); err != nil {
		return err
	}
	if err := c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	hash := common.HashPassword(form.Password)
	member, err := h.db.GetMemberByUsernameAndPassword(c.Request().Context(), db.GetMemberByUsernameAndPasswordParams{
		Username: form.Username,
		Password: string(hash),
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, common.NewErrorJson("failed to login", nil))
	}
	claims := common.NewJwtAuthentication(member.ID, member.Username, member.Password, time.Duration(form.ExpireAfter)*time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(h.tokenSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to create token", err))
	}

	return c.JSON(http.StatusOK, &LoginResponse{
		ID:       member.ID,
		Username: member.Username,
		Role:     string(member.Role),
		Token:    t,
	})
}
