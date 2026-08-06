package services

import (
	"os"

	"github.com/resend/resend-go/v3"
)

type EmailSendService interface {
	SendEmail(emailReq EmailRequest) error
}

type EmailRequest struct {
    From        string
    To          []string
    Subject     string
    HTML        string
    Attachments []*resend.Attachment
}

type resendEmailService struct {
	client *resend.Client
}

func NewEmailSendService() EmailSendService {
	return &resendEmailService{
		client: resend.NewClient(os.Getenv("RESEND_API_KEY")),
	}
}

// Method to send an email using the Resend API
func (e *resendEmailService) SendEmail(emailReq EmailRequest) error {
    params := &resend.SendEmailRequest{
        From:    emailReq.From,
        To:      emailReq.To,
        Subject: emailReq.Subject,
        Html:    emailReq.HTML,
        Attachments: emailReq.Attachments,
    }

    _, err := e.client.Emails.Send(params)
    
    if err != nil {
        return err
    }

	return nil
}