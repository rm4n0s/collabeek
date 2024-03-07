package common

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JwtAuthenticationClaims struct {
	MemberID int32  `json:"memberId"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewJwtAuthentication(memberID int32, name, role string, dur time.Duration) *JwtAuthenticationClaims {
	return &JwtAuthenticationClaims{
		memberID,
		name,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
		},
	}

}

func IsAdminOrModerator(c echo.Context) bool {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return false
	}
	claims, ok := token.Claims.(JwtAuthenticationClaims)
	if !ok {
		return false
	}
	return claims.Role == "admin" || claims.Role == "moderator"
}

func OnlyAdminOrModerators(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !IsAdminOrModerator(c) {
			return NewErrorJson("unauthorized", nil)
		}

		return next(c)
	}
}
