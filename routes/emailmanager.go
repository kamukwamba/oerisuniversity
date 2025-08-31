
package routes

import (
  
    "fmt"
    "net/smtp"
    "os"
    "github.com/jordan-wright/email"
)

// EmailService handles all email operations
type EmailService struct {
    SMTPHost     string
    SMTPPort     int
    SMTPUsername string
    SMTPPassword string
    FromEmail    string
}

// NewEmailService creates a new email service instance
func NewEmailService(host string, port int, username, password, fromEmail string) *EmailService {
    return &EmailService{
        SMTPHost:     host,
        SMTPPort:     port,
        SMTPUsername: username,
        SMTPPassword: password,
        FromEmail:    fromEmail,
    }
}

// SendWelcomeEmail sends a welcome email to new users
func (es *EmailService) SendWelcomeEmail(toEmail, userName string) error {
    subject := "Welcome to Oceries University Colledge Of Metaphysical Sciences"
    
    // Text version
    textBody := fmt.Sprintf(`
        Hello %s,

        Welcome to Oceries University Colledge Of Metaphysical Sciences! We're excited to have you on board.

        If you have any questions, please don't hesitate to contact us.

        Best regards,
        The Team
`, userName)

    // HTML version
    htmlBody := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <style>
                body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
                .container { max-width: 600px; margin: 0 auto; padding: 20px; }
                .header { background: #f8f9fa; padding: 20px; text-align: center; border-radius: 5px; }
                .content { padding: 20px; }
                .footer { margin-top: 20px; padding: 20px; background: #f8f9fa; border-radius: 5px; text-align: center; font-size: 14px; color: #666; }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <h1>Welcome to Oceries University Colledge Of Metaphysical Sciences!</h1>
                </div>
                <div class="content">
                    <p>Hello <strong>%s</strong>,</p>
                    <p>Welcome to Oceries University Colledge Of Metaphysical Sciences! We're excited to have you on board.</p>
                    <p>To log into your account use your email as both your  <span style="font-weight">Username</span> and <span style="font-weight">Password</span>.</p>
                    <p>If you have any questions, please don't hesitate to contact us.</p>
                </div>
                <div class="footer">
                    <p>Best regards,<br>The Team</p>
                </div>
            </div>
        </body>
        </html>
        `, userName)

    return es.sendEmail(toEmail, subject, textBody, htmlBody)
}

// SendPasswordResetEmail sends a password reset email
func (es *EmailService) SendPasswordResetEmail(toEmail, resetToken string) error {
    resetLink := fmt.Sprintf("https://yourapp.com/reset-password?token=%s", resetToken)
    
    subject := "Password Reset Request"
    
    // Text version
    textBody := fmt.Sprintf(`
        You requested a password reset.

        Click the link below to reset your password:
        %s

        This link will expire in 1 hour.

        If you didn't request this reset, please ignore this email.

        Best regards,
        The Team
        `, resetLink)

    // HTML version
    htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #fff3cd; padding: 20px; text-align: center; border-radius: 5px; border: 1px solid #ffeaa7; }
        .content { padding: 20px; }
        .button { display: inline-block; padding: 12px 24px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; margin: 15px 0; }
        .footer { margin-top: 20px; padding: 20px; background: #f8f9fa; border-radius: 5px; text-align: center; font-size: 14px; color: #666; }
        .warning { color: #856404; background: #fff3cd; padding: 10px; border-radius: 3px; border: 1px solid #ffeaa7; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Password Reset Request</h2>
        </div>
        <div class="content">
            <p>You requested a password reset for your account.</p>
            <p>Click the button below to reset your password:</p>
            
            <p style="text-align: center;">
                <a href="%s" class="button">Reset Password</a>
            </p>
            
            <p>Or copy and paste this link in your browser:</p>
            <p><code>%s</code></p>
            
            <div class="warning">
                <p><strong>⚠️ This link will expire in 1 hour.</strong></p>
            </div>
            
            <p>If you didn't request this reset, please ignore this email and your password will remain unchanged.</p>
        </div>
        <div class="footer">
            <p>Best regards,<br>The Team</p>
        </div>
    </div>
</body>
</html>
`, resetLink, resetLink)

    return es.sendEmail(toEmail, subject, textBody, htmlBody)
}

// SendVerificationEmail sends an email verification email
func (es *EmailService) SendVerificationEmail(toEmail, verificationToken string) error {
    verifyLink := fmt.Sprintf("https://yourapp.com/verify-email?token=%s", verificationToken)
    
    subject := "Verify Your Email Address"
    
    textBody := fmt.Sprintf(`
Please verify your email address by clicking the link below:

%s

This verification link will expire in 24 hours.

If you didn't create an account with us, please ignore this email.

Best regards,
The Team
`, verifyLink)

    htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #d4edda; padding: 20px; text-align: center; border-radius: 5px; }
        .button { display: inline-block; padding: 12px 24px; background: #28a745; color: white; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>Verify Your Email Address</h2>
        </div>
        <p>Please verify your email address by clicking the button below:</p>
        <p style="text-align: center;">
            <a href="%s" class="button">Verify Email</a>
        </p>
        <p>Or copy this link to your browser:<br><code>%s</code></p>
        <p>This verification link will expire in 24 hours.</p>
        <p>If you didn't create an account with us, please ignore this email.</p>
        <p>Best regards,<br>The Team</p>
    </div>
</body>
</html>
`, verifyLink, verifyLink)

    return es.sendEmail(toEmail, subject, textBody, htmlBody)
}

// sendEmail is the internal method that actually sends the email
func (es *EmailService) sendEmail(toEmail, subject, textBody, htmlBody string) error {
    e := email.NewEmail()
    e.From = es.FromEmail
    e.To = []string{toEmail}
    e.Subject = subject
    e.Text = []byte(textBody)
    e.HTML = []byte(htmlBody)

    // Add headers for better email deliverability
    e.Headers.Add("X-Mailer", "GoEmailService")
    e.Headers.Add("Precedence", "bulk")
    
    auth := smtp.PlainAuth("", es.SMTPUsername, es.SMTPPassword, es.SMTPHost)
    addr := fmt.Sprintf("%s:%d", es.SMTPHost, es.SMTPPort)

    // Add timeout to prevent hanging
    return e.SendWithTLS(addr, auth, nil)
}

// SendEmail sends a custom email with both text and HTML versions
func (es *EmailService) SendEmail(toEmail, subject, textContent, htmlContent string) error {
    return es.sendEmail(toEmail, subject, textContent, htmlContent)
}

// LoadConfigFromEnv loads email configuration from environment variables
func LoadConfigFromEnv() *EmailService {
    host := getEnv("SMTP_HOST", "smtp.gmail.com")
    port := 587 // default port for SMTP
    username := getEnv("SMTP_USERNAME", "")
    password := getEnv("SMTP_PASSWORD", "")
    fromEmail := getEnv("FROM_EMAIL", "noreply@ocerisumps@gmail.com")

    return NewEmailService(host, port, username, password, fromEmail)
}

// getEnv helper function to get environment variables with defaults
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}




