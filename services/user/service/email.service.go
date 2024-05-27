package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"monorepo/internal/config"
	"monorepo/internal/dto"
	"monorepo/pkg/common"
	"monorepo/pkg/utils"
	"monorepo/services/user/models"
	"text/template"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/oklog/ulid/v2"
	"gopkg.in/gomail.v2"
)

func NewEmailService(
	dialer *gomail.Dialer,
	mailer *gomail.Message,
	fbaClient *auth.Client,
	tbUser common.Repository[models.User, string],
	tbProfile common.Repository[models.Profile, string],
	tbResetPassword common.Repository[models.ResetPassword, string],
) *EmailService {
	service := &EmailService{}
	service.dialer = dialer
	service.mailer = mailer
	service.fbaClient = fbaClient
	service.tables.user = tbUser
	service.tables.profile = tbProfile
	service.tables.resetPassword = tbResetPassword

	return service
}

type EmailService struct {
	dialer    *gomail.Dialer
	mailer    *gomail.Message
	fbaClient *auth.Client
	tables    struct {
		user          common.Repository[models.User, string]
		profile       common.Repository[models.Profile, string]
		resetPassword common.Repository[models.ResetPassword, string]
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

	token := utils.RandAlphanumericString(16)
	rp := &models.ResetPassword{
		ID:         ulid.Make().String(),
		UserID:     user[0].ID,
		ResetToken: token,
		CreatedAt:  time.Now(),
	}
	err = service.tables.resetPassword.Create(ctx, rp)
	if err != nil {
		return fmt.Errorf("%w; %w", ErrRepositoryMutateFail, err)
	}

	sendEmail := dto.Email{
		From:          env.SMTPAuthEmail,
		To:            body.Email,
		Subject:       "Reset Password",
		TemplateEmail: "template/forgot-password.html",
		Body: dto.EmailBody{
			UserName:         name,
			ResetPasswordUrl: fmt.Sprintf("%s?uid=%s&reset-token=%s", env.ResetPasswordUrl, user[0].ID, token),
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

func (service *EmailService) UpdatePassword(ctx context.Context, env *config.Environment, body *dto.RequestUpdatePassword) error {
	user, err := service.tables.user.List(ctx, &common.FilterOptions{
		Sort:   []exp.Expression{goqu.I("id").Desc()},
		Filter: []exp.Expression{goqu.C("id").Eq(body.UserID)},
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

	logReset, err := service.tables.resetPassword.List(ctx, &common.FilterOptions{
		Sort:   []exp.Expression{goqu.I("created_at").Desc()},
		Filter: []exp.Expression{goqu.C("user_id").Eq(body.UserID), goqu.C("is_used").Eq(false)},
		Page:   1,
		Limit:  1,
	})
	if err != nil {
		return err
	}

	if len(logReset) == 0 {
		err = errors.New("request not found")
		return err
	}

	if time.Now().After(logReset[0].CreatedAt.Add(time.Hour * 1)) {
		err = errors.New("token expired")
		return err
	}

	if body.ResetToken != logReset[0].ResetToken {
		err = errors.New("token unknown")
		return err
	}

	// get user by email
	u, err := service.fbaClient.GetUserByEmail(ctx, user[0].Handle)
	if err != nil {
		return err
	}

	// update firebase
	params := (&auth.UserToUpdate{}).
		Password(body.Password)
	_, err = service.fbaClient.UpdateUser(ctx, u.UID, params)
	if err != nil {
		return err
	}

	// update flag is_used
	logReset[0].IsUsed = true
	err = service.tables.resetPassword.Update(ctx, logReset[0].ID, logReset[0])
	if err != nil {
		return err
	}

	return nil
}
