package mail

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v3"
	"github.com/zcl0621/compx576-smart-dairy-system/config"
)

var sendEmail = func(ctx context.Context, client *resend.Client, params *resend.SendEmailRequest) error {
	_, err := client.Emails.SendWithContext(ctx, params)
	return err
}

func SendResetCode(toEmail, code string) error {
	resendConfig := config.Get().Resend
	client := resend.NewClient(resendConfig.APIKey)
	params := &resend.SendEmailRequest{
		From:    resendConfig.From,
		To:      []string{toEmail},
		Subject: "Reset code",
		Text:    fmt.Sprintf("Your code is %s", code),
		Html:    fmt.Sprintf("<p>Your code is <strong>%s</strong></p>", code),
	}

	return sendEmail(context.Background(), client, params)
}
