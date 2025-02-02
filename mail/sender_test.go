package mail_test

import (
	"testing"

	"example.com/simplebank/mail"
	"example.com/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := mail.NewGmailSender(
		config.EmailSenderName,
		config.EmailSenderAddress,
		config.EmailSenderPassword,
	)

	subject := "A test email"
	body := `
		<h1> Hello World</h1>
		<p> This is a test email from FinGo</p>
	`
	to := []string{"workspace.omkarbhostekar@gmail.com"}
	files := []string{"../README.md"}

	err = sender.SendEmail(subject, body, to, nil, nil, files)

	require.NoError(t, err)
}