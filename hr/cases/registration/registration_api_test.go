package registration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rm4n0s/collabeek/common"
	"github.com/rm4n0s/collabeek/hr/db"
	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	email     string
	rs        string
	confirmed bool
	username  string
}

func (m *MockDB) UpdateMemberForRegistration(ctx context.Context, arg db.UpdateMemberForRegistrationParams) error {
	return nil
}

func (m *MockDB) GetMemberByEmail(ctx context.Context, email string) (db.Member, error) {
	if m.email == email {
		member := db.Member{
			ID:                 1,
			Email:              m.email,
			RegistrationSecret: m.rs,
			EmailConfirmed:     m.confirmed,
		}

		return member, nil

	}
	return db.Member{}, sql.ErrNoRows
}

func (m *MockDB) GetMemberByUsername(ctx context.Context, username string) (db.Member, error) {
	if m.username == username {
		member := db.Member{
			ID:                 1,
			Email:              m.email,
			RegistrationSecret: m.rs,
			EmailConfirmed:     m.confirmed,
		}

		return member, nil

	}
	return db.Member{}, sql.ErrNoRows
}

type MockSmtp struct {
	sended bool
}

func (m *MockSmtp) SendEmail(subject string, to []string, msg []byte) error {
	m.sended = true
	return nil
}

func TestRegistrationOnEmptyValueFailure(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	mockDB := &MockDB{}
	mockSmtp := &MockSmtp{}
	input := RegistrationForm{}
	inputJson, err := json.Marshal(input)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewRegistrationHandler(mockDB, mockSmtp)
	if assert.NoError(t, h.RegistrationHandler(c)) {
		errs := []common.ValidatorErrorJson{}
		json.Unmarshal(rec.Body.Bytes(), &errs)
		mp := map[string]string{}
		for _, e := range errs {
			mp[e.Field] = e.Message
		}
		assert.Equal(t, len(mp), 5)
		_, ok := mp["Email"]
		assert.True(t, ok)
		_, ok = mp["Username"]
		assert.True(t, ok)
		_, ok = mp["Fullname"]
		assert.True(t, ok)
		_, ok = mp["RegistrationSecret"]
		assert.True(t, ok)
		_, ok = mp["Password"]
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

}

func TestRegistrationOnValidationFailure(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	t.Run("non existent email", func(t *testing.T) {
		mockDB := &MockDB{}
		mockSmtp := &MockSmtp{}
		rs, err := common.RandomString(60)
		assert.NoError(t, err)
		input := RegistrationForm{
			Email:              "test@gmail.com",
			Username:           "username",
			Password:           "password",
			Fullname:           "full name",
			RegistrationSecret: rs,
		}
		inputJson, err := json.Marshal(input)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := NewRegistrationHandler(mockDB, mockSmtp)
		if assert.NoError(t, h.RegistrationHandler(c)) {
			ej := common.ErrorJson{}
			err = json.Unmarshal(rec.Body.Bytes(), &ej)
			assert.NoError(t, err)
			assert.Equal(t, "failed to find email", ej.Message)
		}
	})

	t.Run("email already confirmed", func(t *testing.T) {
		email := "test@gmail.com"
		mockDB := &MockDB{email: email, confirmed: true}
		mockSmtp := &MockSmtp{}
		rs, err := common.RandomString(60)
		assert.NoError(t, err)
		input := RegistrationForm{
			Email:              email,
			Username:           "username",
			Password:           "password",
			Fullname:           "full name",
			RegistrationSecret: rs,
		}
		inputJson, err := json.Marshal(input)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := NewRegistrationHandler(mockDB, mockSmtp)
		if assert.NoError(t, h.RegistrationHandler(c)) {
			ej := common.ErrorJson{}
			err = json.Unmarshal(rec.Body.Bytes(), &ej)
			assert.NoError(t, err)
			assert.Equal(t, "already registered", ej.Message)
		}
	})

	t.Run("not correct registration secret", func(t *testing.T) {
		email := "test@gmail.com"
		mockDB := &MockDB{email: email, confirmed: false, rs: "aaaaaa"}
		mockSmtp := &MockSmtp{}
		rs, err := common.RandomString(60)
		assert.NoError(t, err)
		input := RegistrationForm{
			Email:              email,
			Username:           "username",
			Password:           "password",
			Fullname:           "full name",
			RegistrationSecret: rs,
		}
		inputJson, err := json.Marshal(input)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := NewRegistrationHandler(mockDB, mockSmtp)
		if assert.NoError(t, h.RegistrationHandler(c)) {
			ej := common.ErrorJson{}
			err = json.Unmarshal(rec.Body.Bytes(), &ej)
			assert.NoError(t, err)
			assert.Equal(t, "incorrect registration secret", ej.Message)
		}
	})

	t.Run("username already taken", func(t *testing.T) {
		rs, err := common.RandomString(60)
		email := "test@gmail.com"
		username := "username"
		mockDB := &MockDB{email: email, confirmed: false, rs: rs, username: username}
		mockSmtp := &MockSmtp{}
		assert.NoError(t, err)
		input := RegistrationForm{
			Email:              email,
			Username:           username,
			Password:           "password",
			Fullname:           "full name",
			RegistrationSecret: rs,
		}
		inputJson, err := json.Marshal(input)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := NewRegistrationHandler(mockDB, mockSmtp)
		if assert.NoError(t, h.RegistrationHandler(c)) {
			ej := common.ErrorJson{}
			err = json.Unmarshal(rec.Body.Bytes(), &ej)
			assert.NoError(t, err)
			assert.Equal(t, "username is used", ej.Message)
		}
	})
}

func TestRegistrationOnSuccess(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	rs, err := common.RandomString(60)
	email := "test@gmail.com"
	username := "username"
	mockDB := &MockDB{email: email, confirmed: false, rs: rs, username: ""}
	mockSmtp := &MockSmtp{}
	assert.NoError(t, err)
	input := RegistrationForm{
		Email:              email,
		Username:           username,
		Password:           "password",
		Fullname:           "full name",
		RegistrationSecret: rs,
	}
	inputJson, err := json.Marshal(input)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewRegistrationHandler(mockDB, mockSmtp)
	if assert.NoError(t, h.RegistrationHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, mockSmtp.sended)
	}
}
