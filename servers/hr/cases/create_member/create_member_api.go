package createmember

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rm4n0s/collabeek/servers/common"
	"github.com/rm4n0s/collabeek/servers/hr/db"
	emailtemplates "github.com/rm4n0s/collabeek/servers/hr/email_templates"
)

func NewCreateMemberHandler(db CreateMemberDB, smtpService common.ISmtpService) *CreateMemberHandler {
	return &CreateMemberHandler{
		db:          db,
		smtpService: smtpService,
	}
}

func (h *CreateMemberHandler) CreateMemberHandler(c echo.Context) error {
	form := new(CreateMemberForm)
	if err := c.Bind(form); err != nil {
		return err
	}
	if err := c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	_, err := h.db.GetMemberByEmail(c.Request().Context(), form.Email)
	if err == nil {
		return c.JSON(http.StatusBadRequest, common.NewErrorJson("member exists", nil))
	} else {
		if err != sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to search member", err))
		}
	}

	rs, err := common.RandomString(60)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to create random", err))
	}

	obj := db.CreateMemberParams{
		Email:              form.Email,
		Role:               db.Roles(form.Role),
		RegistrationSecret: rs,
	}
	member, err := h.db.CreateMember(c.Request().Context(), obj)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to create member", err))
	}

	var email bytes.Buffer
	t, err := template.New("registration").Parse(emailtemplates.RegistrationLetter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to parse email tempalte", err))
	}
	err = t.Execute(&email, &emailtemplates.RegistrationLetterInput{RegistrationCode: rs})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to create email", err))
	}
	err = h.smtpService.SendEmail("Registration code", []string{form.Email}, email.Bytes())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to send email", err))
	}

	cm := CreateMemberResponse{
		ID:        member.ID,
		Email:     member.Email,
		Role:      string(member.Role),
		CreatedAt: member.CreatedAt,
	}
	return c.JSON(http.StatusCreated, cm)
}
