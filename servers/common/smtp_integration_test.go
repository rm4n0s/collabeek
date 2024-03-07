package common

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mailhogListJson struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Start int `json:"start"`
	Items []struct {
		ID   string `json:"ID"`
		From struct {
			Relays  any    `json:"Relays"`
			Mailbox string `json:"Mailbox"`
			Domain  string `json:"Domain"`
			Params  string `json:"Params"`
		} `json:"From"`
		To []struct {
			Relays  any    `json:"Relays"`
			Mailbox string `json:"Mailbox"`
			Domain  string `json:"Domain"`
			Params  string `json:"Params"`
		} `json:"To"`
		Content struct {
			Headers struct {
				ContentType []string `json:"Content-Type"`
				MIMEVersion []string `json:"MIME-version"`
				MessageID   []string `json:"Message-ID"`
				Received    []string `json:"Received"`
				ReturnPath  []string `json:"Return-Path"`
				Subject     []string `json:"Subject"`
			} `json:"Headers"`
			Body string `json:"Body"`
			Size int    `json:"Size"`
			Mime any    `json:"MIME"`
		} `json:"Content"`
		Created time.Time `json:"Created"`
		Mime    any       `json:"MIME"`
		Raw     struct {
			From string   `json:"From"`
			To   []string `json:"To"`
			Data string   `json:"Data"`
			Helo string   `json:"Helo"`
		} `json:"Raw"`
	} `json:"items"`
}

func TestSmtpServiceWithMailhog(t *testing.T) {
	ss := NewSmtpService("from@test.com", "test", "test", "127.0.0.1", 1025)
	rstr, err := RandomString(10)
	assert.NoError(t, err)
	subject := "test email " + rstr
	err = ss.SendEmail(subject, []string{"to@to.com"}, []byte("message test"))
	assert.NoError(t, err)
	client := &http.Client{}
	getListReq, err := http.NewRequest("GET", "http://localhost:8025/api/v2/messages?limit=50", nil)
	getListReq.SetBasicAuth("test", "test")
	assert.NoError(t, err)
	resp, err := client.Do(getListReq)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	mlj := mailhogListJson{}
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(b, &mlj)
	assert.NoError(t, err)
	exists := false
	for _, v := range mlj.Items {
		itemSubject := strings.TrimSpace(strings.Join(v.Content.Headers.Subject, ""))
		if itemSubject == strings.TrimSpace(subject) {
			exists = true
			break
		}
	}
	assert.True(t, exists)
}
