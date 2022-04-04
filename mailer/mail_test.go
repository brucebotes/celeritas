package mailer

import (
	"errors"
	"testing"
)

func TestMail_SendSMTPMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
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
		FromName:    "Joe",
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

	msg.To = "not_an_email_address"
	mailer.Jobs <- msg
	res = <-mailer.Results
	if res.Error == nil {
		t.Error(errors.New("mail should have failed with an invalid email address but it did not"))
	}
}

func TestMail_SendUsingAPI(t *testing.T) {
	msg := Message{
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.API = "unknown"
	mailer.APIKey = "abc123"
	mailer.APIUrl = "https://www.fake.com"

	err := mailer.SendUsingApi(msg, "unkown")
	if err == nil {
		t.Error(errors.New("expected an error but did not get one"))
	}

	// reset values for tests following this one
	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""
}

func TestMail_buildHTMLMessage(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
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
		FromName:    "Joe",
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

func TestMail_send(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	err := mailer.Send(msg)
	if err != nil {
		t.Error(err)
	}

	mailer.API = "unknown"
	mailer.APIKey = "abc123"
	mailer.APIUrl = "https://www.fake.com"

	err = mailer.Send(msg)
	if err == nil {
		t.Error(errors.New("Send - expected an error but did not get one"))
	}
	// reset values for tests following this one
	mailer.API = ""
	mailer.APIKey = ""
	mailer.APIUrl = ""
}

func TestMail_ChooseAPI(t *testing.T) {
	msg := Message{
		From:        "me@here.com",
		FromName:    "Joe",
		To:          "you@there.com",
		Subject:     "test",
		Template:    "test",
		Attachments: []string{"./testdata/mail/test.html.tmpl"},
	}

	mailer.API = "unknown"
	mailer.APIKey = "abc123"
	mailer.APIUrl = "https://www.fake.com"

	err := mailer.ChooseAPI(msg)
	if err == nil {
		t.Error(errors.New("ChooseAPI - expected an error but did not get one"))
	}
}
