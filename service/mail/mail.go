package mail

import (
	"carb-backend/cmd/db"
	"fmt"

	"github.com/wneessen/go-mail"
)

func VerificationTemplate(username, url string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Verify your e-mail</title>
		</head>

		<body style="
			margin: 0;
			padding: 0;
			background-color: #f4f4f5;
			font-family: Arial, Helvetica, sans-serif;
			color: #18181b;
		">

			<div style="
				max-width: 600px;
				margin: 40px auto;
				padding: 0 20px;
			">

				<div style="
					background: #ffffff;
					border-radius: 12px;
					padding: 40px;
					border: 1px solid #e4e4e7;
				">

					<h1 style="
						margin: 0 0 24px;
						font-size: 28px;
					">
						Welcome, %s!
					</h1>

					<p style="
						font-size: 16px;
						line-height: 1.6;
						color: #52525b;
					">
						Thanks for creating an account.
						Before you can continue, please verify your e-mail address.
					</p>

					<div style="text-align: center; margin: 32px 0;">
						<a href="%s"
						style="
							display: inline-block;
							padding: 14px 28px;
							background-color: #18181b;
							color: #ffffff;
							text-decoration: none;
							border-radius: 8px;
							font-size: 16px;
							font-weight: 600;
						">
							Verify my email
						</a>
					</div>

					<p style="
						font-size: 14px;
						line-height: 1.6;
						color: #71717a;
					">
						This verification link will expire in 24 hours.
					</p>

					<hr style="
						border: none;
						border-top: 1px solid #e4e4e7;
						margin: 32px 0;
					">

					<p style="
						font-size: 13px;
						line-height: 1.5;
						color: #a1a1aa;
					">
						If you didn't create an account, you can safely ignore this email.
					</p>

				</div>

				<p style="
					text-align: center;
					margin-top: 24px;
					font-size: 12px;
					color: #a1a1aa;
				">
					© 2026 Go Backend. All rights reserved.
				</p>

			</div>

		</body>
	</html>
	`, username, url)
}

func SendEmail(to, subject, body string, debug bool) error {
	host := db.Environement.SMTPHost
	port := db.Environement.SMTPPort
	from := db.Environement.SMTPFrom

	m := mail.NewMsg()
	m.From(from)
	m.To(to)
	m.Subject(subject)

	m.SetBodyString(mail.TypeTextHTML, body)

	var client *mail.Client
	var err error

	if debug {
		client, err = mail.NewClient(
			host,
			mail.WithPort(port),
			mail.WithTLSPolicy(mail.NoTLS),
		)
	} else {
		client, err = mail.NewClient(
			host,
			mail.WithPort(port),
			mail.WithTLSPolicy(mail.TLSMandatory),
		)
	}

	if err != nil {
		return err
	}

	return client.DialAndSend(m)
}
