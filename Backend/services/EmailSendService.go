package services

import (
	"os"

	"github.com/resend/resend-go/v3"
)

type EmailSendService interface {
	SendEmail() error
}

type emailSendService struct {
	client *resend.Client
}

func NewEmailSendService() EmailSendService {
	return &emailSendService{
		client: resend.NewClient(os.Getenv("RESEND_API_KEY")),
	}
}

// Method to send an email using the Resend API
func (e *emailSendService) SendEmail() error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

    params := &resend.SendEmailRequest{
        From:    "ScheduleFlow <noreply@darienmiller.com>",
        To:      []string{"darienm931@gmail.com"},
        Subject: "Hello World",
        Html:    "<p>Congrats on sending your <strong>first email</strong>!</p>",
    }

    _, err := client.Emails.Send(params)
    
    if err != nil {
        return err
    }

	return nil
}