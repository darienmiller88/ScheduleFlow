package services

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

type EmailSendService interface {
	SendEmail(emailReq EmailRequest) error
	SendVerificationEmail(toEmail string, firstName string, verificationCode string) error
}

type EmailRequest struct {
	From     string
	To       []string
	CC       []string
	ReplyTo  string
	Subject  string
	HTML     string
	Filepath string
}

type resendEmailService struct {
	client *resend.Client
}


func NewEmailSendService() EmailSendService {
	return &resendEmailService{
		client: resend.NewClient(os.Getenv("RESEND_API_KEY")),
	}
}

// SendVerificationEmail implements [EmailSendService].
func (e *resendEmailService) SendVerificationEmail(toEmail string, firstName string, verificationCode string) error {
	_, err := e.client.Emails.Send(&resend.SendEmailRequest{
		From:    "ScheduleFlow <noreply@darienmiller.com>",
		To:      []string{toEmail},
		Subject: "Verify your email address",
		Html:    fmt.Sprintf(`
		<div style="margin:0; padding:40px 20px; background-color:#f4f6f8; font-family:Arial, Helvetica, sans-serif;">
			<div style="max-width:600px; margin:0 auto; background-color:#ffffff; border:1px solid #e1e5e9; border-radius:8px; overflow:hidden;">

				<!-- Header -->
				<div style="padding:24px 30px; background-color:#1f2937;">
					<h1 style="margin:0; color:#ffffff; font-size:24px; font-weight:600;">
						ScheduleFlow
					</h1>
				</div>

				<!-- Content -->
				<div style="padding:30px;">
					<h2 style="margin:0 0 16px 0; color:#111827; font-size:20px;">
						Verify your email address
					</h2>

					<p style="margin:0 0 16px 0; color:#4b5563; font-size:15px; line-height:1.6;">
						Hello %s,
					</p>

					<p style="margin:0 0 24px 0; color:#4b5563; font-size:15px; line-height:1.6;">
						Thank you for signing up for ScheduleFlow. 
						Enter the verification code below to confirm your email address.
					</p>

					<!-- Verification Code -->
					<div style="margin:0 0 24px 0; padding:20px; background-color:#f3f4f6; border:1px solid #e5e7eb; border-radius:8px; text-align:center;">
						<p style="margin:0 0 8px 0; color:#6b7280; font-size:12px; text-transform:uppercase; letter-spacing:1px;">
							Verification Code
						</p>

						<p style="margin:0; color:#111827; font-size:32px; font-weight:700; letter-spacing:8px;">
							%s
						</p>
					</div>

					<p style="margin:0 0 12px 0; color:#6b7280; font-size:13px; line-height:1.5;">
						This code will expire in <strong>15 minutes</strong>.
					</p>

					<p style="margin:0; color:#6b7280; font-size:13px; line-height:1.5;">
						If you did not create a ScheduleFlow account, you can safely ignore this email.
					</p>
				</div>

				<!-- Footer -->
				<div style="padding:20px 30px; border-top:1px solid #e5e7eb;">
					<p style="margin:0; color:#9ca3af; font-size:12px; text-align:center;">
						This is an automated message from ScheduleFlow.
					</p>
				</div>

			</div>
		</div>
	`, firstName, verificationCode),
	})

	if err != nil {
		return err
	}

	return nil
}

// Method to send an email using the Resend API
func (e *resendEmailService) SendEmail(emailReq EmailRequest) error {
	_, err := e.client.Emails.Send(&resend.SendEmailRequest{
		From:    emailReq.From,
		To:      emailReq.To,
		Cc:      emailReq.CC,
		ReplyTo: emailReq.ReplyTo,
		Subject: emailReq.Subject,
		Html:    emailReq.HTML,
		Attachments: []*resend.Attachment{
			{
				Filename: emailReq.Filepath,
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
