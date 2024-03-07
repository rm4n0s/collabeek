package isauthenticated

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rm4n0s/collabeek/common"
)

func NewIsAuthenticatedHandler(tokenSecret string) *IsAuthenticatedHandler {
	return &IsAuthenticatedHandler{
		tokenSecret: []byte(tokenSecret),
	}
}

func (h *IsAuthenticatedHandler) IsAuthenticatedHandler(c echo.Context) error {
	form := new(IsAuthenticatedForm)
	if err := c.Bind(form); err != nil {
		return err
	}
	if err := c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	token, err := jwt.Parse(form.Token, func(token *jwt.Token) (interface{}, error) {
		return h.tokenSecret, nil
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewErrorJson("failed to parse token", err))
	}

	if !token.Valid {
		return c.JSON(http.StatusUnauthorized, common.NewErrorJson("unauthorized", nil))
	}

	return c.NoContent(http.StatusOK)
}
