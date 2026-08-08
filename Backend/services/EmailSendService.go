package services

import (
	"os"

	"github.com/resend/resend-go/v3"
)

type EmailSendService interface {
	SendEmail(emailReq resend.SendEmailRequest) error
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
func (e *resendEmailService) SendEmail(emailReq resend.SendEmailRequest) error {
    _, err := e.client.Emails.Send(&emailReq)
    
    if err != nil {
        return err
    }

	return nil
}