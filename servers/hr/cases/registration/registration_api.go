package registration

import (
	"bytes"
	"database/sql"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/rm4n0s/collabeek/servers/common"
	"github.com/rm4n0s/collabeek/servers/hr/db"
	emailtemplates "github.com/rm4n0s/collabeek/servers/hr/email_templates"
)

func NewRegistrationHandler(db RegistrationDB, smtpService common.ISmtpService) *RegistrationHandler {
	return &RegistrationHandler{
		db:          db,
		smtpService: smtpService,
	}
}

func (r *RegistrationHandler) RegistrationHandler(c echo.Context) error {
	form := new(RegistrationForm)
	if err := c.Bind(form); err != nil {
		return err
	}
	if err := c.Validate(form); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	member, err := r.db.GetMemberByEmail(c.Request().Context(), form.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusBadRequest, common.NewErrorJson("failed to find email", nil))
		}
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("problem with database", err))
	}
	if member.EmailConfirmed {
		return c.JSON(http.StatusBadRequest, common.NewErrorJson("already registered", nil))
	}
	if member.RegistrationSecret != form.RegistrationSecret {
		return c.JSON(http.StatusBadRequest, common.NewErrorJson("incorrect registration secret", nil))
	}

	_, err = r.db.GetMemberByUsername(c.Request().Context(), form.Username)
	if err == nil {
		return c.JSON(http.StatusBadRequest, common.NewErrorJson("username is used", nil))
	} else {
		if err != sql.ErrNoRows {
			return c.JSON(http.StatusInternalServerError, common.NewErrorJson("problem with database", err))
		}
	}

	hash := common.HashPassword(form.Password)

	err = r.db.UpdateMemberForRegistration(c.Request().Context(), db.UpdateMemberForRegistrationParams{
		ID:                 member.ID,
		Username:           form.Username,
		Password:           hash,
		Fullname:           form.Fullname,
		EmailConfirmed:     true,
		RegistrationSecret: "",
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to update member", err))
	}

	var email bytes.Buffer
	t, err := template.New("registration").Parse(emailtemplates.SuccessfulRegistrationLetter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to parse email tempalte", err))
	}
	err = t.Execute(&email, &emailtemplates.SuccessfulRegistrationLetterInput{Fullname: form.Fullname})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to create email", err))
	}

	err = r.smtpService.SendEmail("Registration succeed", []string{form.Email}, email.Bytes())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewErrorJson("failed to send email", err))
	}
	return c.JSON(http.StatusOK, nil)
}
