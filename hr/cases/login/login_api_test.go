package login

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rm4n0s/collabeek/common"
	"github.com/rm4n0s/collabeek/hr/db"
	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	username string
	password string
	role     string
}

func (m *MockDB) GetMemberByUsernameAndPassword(ctx context.Context, args db.GetMemberByUsernameAndPasswordParams) (db.Member, error) {
	fmt.Println(args)
	member := db.Member{
		ID:       1,
		Username: m.username,
		Password: m.password,
		Role:     db.Roles(m.role),
	}
	if args.Username != m.username {
		return db.Member{}, errors.New("not found")
	}
	hash := common.HashPassword(m.password)
	if hash != args.Password {
		return db.Member{}, errors.New("not found")
	}
	return member, nil
}

func TestLoginSuccess(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	input := LoginForm{
		Username:    "user",
		Password:    "password",
		ExpireAfter: time.Second,
	}
	inputJson, err := json.Marshal(&input)
	assert.NoError(t, err)
	mockDB := &MockDB{username: "user", password: "password"}
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewLoginHandler(mockDB, "secret")
	if assert.NoError(t, h.LoginHandler(c)) {
		resp := LoginResponse{}
		json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, resp.Username, mockDB.username)
		assert.Equal(t, resp.ID, int32(1))
		assert.Equal(t, resp.Role, mockDB.role)
		assert.NotEmpty(t, resp.Token)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	successJson := "{\"message\":\"failed to login\",\"error\":\"\"}"
	input := LoginForm{
		Username:    "username",
		Password:    "password",
		ExpireAfter: time.Second,
	}
	inputJson, err := json.Marshal(&input)
	assert.NoError(t, err)
	mockDB := &MockDB{username: "user", password: "password"}
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewLoginHandler(mockDB, "secret")

	if assert.NoError(t, h.LoginHandler(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, successJson, strings.TrimSpace(rec.Body.String()))
	}
}

func TestLoginEmptyForm(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()

	input := LoginForm{
		Username:    "",
		Password:    "",
		ExpireAfter: -1,
	}
	inputJson, err := json.Marshal(&input)
	assert.NoError(t, err)
	mockDB := &MockDB{username: "user", password: "password"}
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewLoginHandler(mockDB, "secret")

	if assert.NoError(t, h.LoginHandler(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		errs := []common.ValidatorErrorJson{}
		json.Unmarshal(rec.Body.Bytes(), &errs)
		mp := map[string]string{}
		for _, e := range errs {
			mp[e.Field] = e.Message
		}
		assert.Equal(t, len(mp), 3)
		_, ok := mp["Username"]
		assert.True(t, ok)
		_, ok = mp["Password"]
		assert.True(t, ok)
		_, ok = mp["ExpireAfter"]
		assert.True(t, ok)
	}
}
