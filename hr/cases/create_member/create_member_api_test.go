package createmember

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rm4n0s/collabeek/common"
	"github.com/rm4n0s/collabeek/hr/db"
	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	email string
	role  string
}

func (m *MockDB) CreateMember(ctx context.Context, args db.CreateMemberParams) (db.Member, error) {
	member := db.Member{
		ID:    1,
		Email: args.Email,
		Role:  db.Roles(args.Role),
	}

	return member, nil
}

func (m *MockDB) GetMemberByEmail(ctx context.Context, email string) (db.Member, error) {
	if m.email == email {
		member := db.Member{
			ID:    1,
			Email: m.email,
			Role:  db.Roles(m.role),
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

func TestCreateMemberOnEmptyValuesFailure(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	mockDB := &MockDB{role: "", email: ""}
	mockSmtp := &MockSmtp{}
	input := CreateMemberForm{}
	inputJson, err := json.Marshal(input)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewCreateMemberHandler(mockDB, mockSmtp)
	if assert.NoError(t, h.CreateMemberHandler(c)) {
		errs := []common.ValidatorErrorJson{}
		json.Unmarshal(rec.Body.Bytes(), &errs)
		mp := map[string]string{}
		for _, e := range errs {
			mp[e.Field] = e.Message
		}
		assert.Equal(t, len(mp), 2)
		_, ok := mp["Email"]
		assert.True(t, ok)
		_, ok = mp["Role"]
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreateMemberOnWrongEmailAndRoleFailure(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	mockDB := &MockDB{role: "", email: ""}
	input := CreateMemberForm{Role: "king", Email: "lalla$gmail.com"}
	mockSmtp := &MockSmtp{}
	inputJson, err := json.Marshal(input)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewCreateMemberHandler(mockDB, mockSmtp)
	if assert.NoError(t, h.CreateMemberHandler(c)) {
		errs := []common.ValidatorErrorJson{}
		json.Unmarshal(rec.Body.Bytes(), &errs)
		mp := map[string]string{}
		for _, e := range errs {
			mp[e.Field] = e.Message
		}
		assert.Equal(t, len(mp), 2)
		_, ok := mp["Email"]
		assert.True(t, ok)
		_, ok = mp["Role"]
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreateMemberOnExistedEmailFailure(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	mockDB := &MockDB{role: "member", email: "lala@gmail.com"}
	mockSmtp := &MockSmtp{}
	input := CreateMemberForm{Role: "member", Email: "lala@gmail.com"}
	inputJson, err := json.Marshal(input)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewCreateMemberHandler(mockDB, mockSmtp)
	failureMsg := `{"message":"member exists","error":""}`
	if assert.NoError(t, h.CreateMemberHandler(c)) {
		assert.Equal(t, failureMsg, strings.TrimSpace(rec.Body.String()))
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestCreateMemberOnSuccess(t *testing.T) {
	e := echo.New()
	e.Validator = common.NewValidator()
	mockDB := &MockDB{role: "", email: ""}
	mockSmtp := &MockSmtp{}
	input := CreateMemberForm{Role: "member", Email: "lala@gmail.com"}
	inputJson, err := json.Marshal(input)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := NewCreateMemberHandler(mockDB, mockSmtp)
	successMsg := "{\"ID\":1,\"Email\":\"lala@gmail.com\",\"Role\":\"member\",\"CreatedAt\":\"0001-01-01T00:00:00Z\"}"
	if assert.NoError(t, h.CreateMemberHandler(c)) {
		assert.Equal(t, successMsg, strings.TrimSpace(rec.Body.String()))
		assert.True(t, mockSmtp.sended)
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
