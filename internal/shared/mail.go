package shared

import (
	"net/mail"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type SendEmailArgs struct {
	ToEmail string
	Subject string
	Text    string
}

func SendEmail(app *pocketbase.PocketBase, args SendEmailArgs) error {
	appMeta := app.Settings().Meta

	message := &mailer.Message{
		From: mail.Address{
			Address: appMeta.SenderAddress,
			Name:    appMeta.SenderName,
		},
		To:      []mail.Address{{Address: args.ToEmail}},
		Subject: args.Subject,
		Text:    args.Text,
	}

	return app.NewMailClient().Send(message)
}
