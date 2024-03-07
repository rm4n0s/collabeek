package isauthenticated

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rm4n0s/collabeek/servers/common"
	"github.com/stretchr/testify/assert"
)

func TestIsAuthenticatedSuccess(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	tokenSecret := "secret"
	claims := common.NewJwtAuthentication(1, "username", "password", time.Second)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := token.SignedString([]byte(tokenSecret))
	assert.NoError(t, err)

	input := IsAuthenticatedForm{Token: tok}
	inputJson, err := json.Marshal(input)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewIsAuthenticatedHandler(tokenSecret)
	if assert.NoError(t, h.IsAuthenticatedHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestIsAuthenticatedOnTimeFailure(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	tokenSecret := "secret"
	claims := common.NewJwtAuthentication(1, "username", "password", time.Second)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := token.SignedString([]byte(tokenSecret))
	assert.NoError(t, err)

	time.Sleep(time.Second)
	input := IsAuthenticatedForm{Token: tok}
	inputJson, err := json.Marshal(input)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	failureMsg := `{"message":"failed to parse token","error":"token has invalid claims: token is expired"}`
	h := NewIsAuthenticatedHandler(tokenSecret)
	if assert.NoError(t, h.IsAuthenticatedHandler(c)) {
		assert.Equal(t, failureMsg, strings.TrimSpace(rec.Body.String()))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestIsAuthenticatedEmptyToken(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	tokenSecret := "secret"
	input := IsAuthenticatedForm{Token: ""}
	inputJson, err := json.Marshal(input)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewIsAuthenticatedHandler(tokenSecret)
	if assert.NoError(t, h.IsAuthenticatedHandler(c)) {
		errs := []common.ValidatorErrorJson{}
		json.Unmarshal(rec.Body.Bytes(), &errs)
		mp := map[string]string{}
		for _, e := range errs {
			mp[e.Field] = e.Message
		}
		assert.Equal(t, len(mp), 1)
		_, ok := mp["Token"]
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}
