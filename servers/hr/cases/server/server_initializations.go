package server

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"text/template"

	"github.com/rm4n0s/collabeek/servers/common"
	"github.com/rm4n0s/collabeek/servers/hr/db"
	emailtemplates "github.com/rm4n0s/collabeek/servers/hr/email_templates"
)

func initializeAdmin(ctx context.Context, idb InitializationDB, smtp common.ISmtpService, adminEmail string) {
	member, err := idb.GetMemberByEmail(ctx, adminEmail)
	if err == nil {
		log.Printf("Admin is: %s", adminEmail)
		log.Printf("Admin is registered: %t", member.EmailConfirmed)
		return
	} else {
		if err != sql.ErrNoRows {
			log.Fatal(fmt.Errorf("failed to search if admin's email exists: %w", err))
		}
	}
	rs, err := common.RandomString(60)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to generate random string on admin's initialization: %w", err))
	}
	member, err = idb.CreateMember(ctx, db.CreateMemberParams{
		Email:              adminEmail,
		Role:               db.RolesAdmin,
		RegistrationSecret: rs,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create admin's membership: %w", err))
	}

	var email bytes.Buffer
	t, err := template.New("registration").Parse(emailtemplates.RegistrationLetter)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse email tempalte in initializeAdmin: %w", err))
	}
	err = t.Execute(&email, &emailtemplates.RegistrationLetterInput{RegistrationCode: rs})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create email in initializeAdmin: %w", err))
	}
	err = smtp.SendEmail("Registration code", []string{adminEmail}, email.Bytes())
	if err != nil {
		log.Fatal(fmt.Errorf("failed to send email in initializeAdmin: %w", err))
	}
	log.Printf("Admin is: %s", adminEmail)
	log.Printf("Admin is registered: %t", member.EmailConfirmed)
}
