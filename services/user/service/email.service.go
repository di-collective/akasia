package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"monorepo/internal/config"
	"monorepo/internal/dto"
	"monorepo/pkg/common"
	"monorepo/services/user/models"
	"text/template"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"gopkg.in/gomail.v2"
)

func NewEmailService(
	dialer *gomail.Dialer,
	mailer *gomail.Message,
	tbUser common.Repository[models.User, string],
	tbProfile common.Repository[models.Profile, string],
) *EmailService {
	service := &EmailService{}
	service.dialer = dialer
	service.mailer = mailer
	service.tables.user = tbUser
	service.tables.profile = tbProfile

	return service
}

type EmailService struct {
	dialer *gomail.Dialer
	mailer *gomail.Message
	tables struct {
		user    common.Repository[models.User, string]
		profile common.Repository[models.Profile, string]
	}
}

func (service *EmailService) ResetPassword(ctx context.Context, env *config.Environment, body *dto.RequestForgotPassword) error {
	user, err := service.tables.user.List(ctx, &common.FilterOptions{
		Sort:   []exp.Expression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("handle").Eq(body.Email)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return err
	}

	if len(user) == 0 {
		err = errors.New("user does not exist")
		return err
	}

	existing, err := service.tables.profile.List(ctx, &common.FilterOptions{
		Sort:   []exp.Expression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("user_id").Eq(user[0].ID)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return err
	}

	if len(existing) == 0 {
		err = errors.New("profile already exist")
		return err
	}

	name := fmt.Sprintf("%s %s", existing[0].FirstName, existing[0].LastName)

	sendEmail := dto.Email{
		From:          env.SMTPAuthEmail,
		To:            body.Email,
		Subject:       "Reset Password",
		TemplateEmail: "template/forgot-password.html",
		Body: dto.EmailBody{
			UserName:         name,
			ResetPasswordUrl: "#",
			CsUrl:            env.CsUrl,
		},
	}

	template, err := service.ParseTemplate(sendEmail.TemplateEmail, sendEmail.Body)
	if err != nil {
		return err
	}

	service.mailer.SetHeader("From", sendEmail.From)
	service.mailer.SetAddressHeader("Cc", sendEmail.To, sendEmail.Subject)
	service.mailer.SetHeader("To", sendEmail.To)
	service.mailer.SetHeader("Subject", sendEmail.Subject)
	service.mailer.SetBody("text/html", template)
	return service.dialer.DialAndSend(service.mailer)
}

func (service *EmailService) ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
