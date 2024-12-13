package mail

import (
	"testing"

	"github.com/bensmile/wekamakuta/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {

	config, err := util.LoadConfig("..")

	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test Email"
	content := `
	
	<h1> Hello World </p>
	`

	to := []string{"benjkafirongo@gmail.com"}

	attachments := []string{"../notes.txt"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachments)

	require.NoError(t, err)
}
