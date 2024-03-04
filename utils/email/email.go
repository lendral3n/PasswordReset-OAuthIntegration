package email

import (
	"bytes"
	"crypto/tls"
	"emailnotifl3n/app/config"
	"emailnotifl3n/features/user"
	"text/template"

	"github.com/k3a/html2text"
	"gopkg.in/gomail.v2"
)

type emailData struct {
	URL     string
	Name    string
	Subject string
}

type emailService struct {
	url      string
	from     string
	password string
	host     string
	port     int
	user     string
}

type EmailInterface interface {
	SendResetPasswordLink(user *user.Core, token string) error
	SendVerificationLink(user *user.Core, token string) error
	SendCodeReset(user *user.Core, code string) error 
}

func New() EmailInterface {
	cfg := config.InitConfig()
	emailCfg := cfg
	return &emailService{
		url:      emailCfg.PASSWD_URL,
		from:     emailCfg.EMAIL_FROM,
		user:     emailCfg.SMTP_USER,
		password: emailCfg.SMTP_PASS,
		host:     emailCfg.SMTP_HOST,
		port:     emailCfg.SMTP_PORT,
	}
}

func (e *emailService) SendResetPasswordLink(user *user.Core, token string) error {
	t, err := template.ParseGlob("utils/templates/resetpasswordlink.html")
	if err != nil {
		return err
	}
	
	to := user.Email
	data := &emailData{
		URL:     e.url + "/reset-password?token=" + token,
		Name:    user.Name,
		Subject: "Reset Password",
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(e.host, e.port, e.user, e.password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (e *emailService) SendVerificationLink(user *user.Core, token string) error {
	t, err := template.ParseGlob("utils/templates/verifiedlink.html")
	if err != nil {
		return err
	}
	
	to := user.Email
	data := &emailData{
		URL:     e.url + "/verification?token=" + token,
		Name:    user.Name,
		Subject: "Email Verification",
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(e.host, e.port, e.user, e.password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func (e *emailService) SendCodeReset(user *user.Core, code string) error {
	t, err := template.ParseGlob("utils/templates/resetpasswordcode.html")
	if err != nil {
		return err
	}
	
	data := &emailData{
		URL:     code,
		Name:    user.Name,
		Subject: "Reset Password Code",
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", e.from)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	d := gomail.NewDialer(e.host, e.port, e.user, e.password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

