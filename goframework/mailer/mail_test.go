package mailer

import (
	"errors"
	"testing"
)

func TestMail_SendSMTPMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "jeje",
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	err := mailer.SendSMTPMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_SendUsingChan(t *testing.T) {

	msg := Message{
		From:        "me@here.com",
		FromName:    "jeje",
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.Jobs <- msg

	res := <-mailer.Results
	if res.Error != nil {
		t.Error(errors.New("failed to send over channel"))
	}

	// test not an email address
	msg.To = "not an email address"
	mailer.Jobs <- msg
	res = <-mailer.Results
	if res.Error == nil {
		t.Error(errors.New("sneding an email without a valid email should return an err"))
	}
}

// can not test the api, so we just test the logic before to send the message
// means we test err in case of incorrect API, apiKeys
func TestMail_SendUsingAPI(t *testing.T) {

	// From & FromNane will be populated automaticaly from our code
	// FromName from .env file
	// From where does it have a default value ????
	msg := Message{
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	// test wrong api
	mailer.API = "unknown"
	mailer.APIKeys = "123abc"
	mailer.APIUrl = "https://www.kake.com"
	err := mailer.SendUsingAPI(msg, "unknow")
	if err == nil {
		t.Error(err)
	}

	// set the mailer data has it was before
	mailer.API = ""
	mailer.APIKeys = ""
	mailer.APIUrl = ""
}

// already tested in func above, but not individually, so lets do it
func TestMail_buildHTMLMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "jeje",
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	_, err := mailer.buildHTMLMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_buildPlainTextMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "jeje",
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	_, err := mailer.buildPlainTextMessage(msg)
	if err != nil {
		t.Error(err)
	}
}

func TestMail_sent(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "jeje",
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}
	err := mailer.Send(msg)
	if err != nil {
		t.Error(err)
	}

	// same as before we test Send() with invalid credentials
	mailer.API = "unknown"
	mailer.APIKeys = "123abc"
	mailer.APIUrl = "https://www.kake.com"
	err = mailer.Send(msg)
	if err == nil {
		t.Error("mailer.send() do not return an err when it should")
	}

	// set the mailer data has it was before
	mailer.API = ""
	mailer.APIKeys = ""
	mailer.APIUrl = ""
}

func TestMail_ChooseAPI(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "jeje",
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}
	// everythings which is not sendgrid, mailgun, smtp or the other one
	mailer.API = "unknown"
	err := mailer.ChooseAPI(msg)
	if err == nil {
		t.Error(err)
	}
}
