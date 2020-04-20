package model

import (
	"database/sql"
	"time"

	"github.com/lvl484/user-manager/logger"

	"github.com/pkg/errors"

	"github.com/lvl484/user-manager/config"

	"github.com/lvl484/user-manager/server/mail"
)

type Verification struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

const (
	queryAddActivationCode          = `INSERT INTO email_codes(id, user_name, email, verification_code, created_at) VALUES ($1,$2,$3,$4,$5)`
	queryDeleteActivationCode       = `DELETE FROM email_codes WHERE user_name=$1`
	querySelectVerificationCodeTime = `SELECT verification_code, created_at FROM email_codes WHERE user_name=$1`
)

// AddActivationCode adds new activation code for user to database
func (ur *UsersRepo) AddActivationCode(user *User) error {
	emailConfig := SetupEmailComponents(user.Email)

	verificationCode := mail.GenerateVerificationCode()

	emailConfig.SendMail(user.Username, verificationCode)

	_, err := ur.db.Exec(queryAddActivationCode, user.ID, user.Username, user.Email, verificationCode, time.Now())
	return err
}

// DeleteVerificationCode deletes activation code for user from database
func (ur *UsersRepo) DeleteVerificationCode(login string) error {
	_, err := ur.db.Exec(queryDeleteActivationCode, login)

	return err
}

// GetVerificationCodeTime gets verification code and time when it was code was created from database
func (ur *UsersRepo) GetVerificationCodeTime(login string) (*time.Time, string, error) {
	var verifyCodeTime *time.Time
	var verificationCode string

	err := ur.db.QueryRow(querySelectVerificationCodeTime, login).Scan(&verificationCode, &verifyCodeTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.Wrap(err, msgUserDidNotExist)
		}
		return nil, "", err
	}

	return verifyCodeTime, verificationCode, nil
}

// SetupEmailComponents setups email components
func SetupEmailComponents(email string) *mail.EmailInfo {
	cfg, err := config.NewConfig()
	if err != nil {
		logger.LogUM.Infof("SetupEmailComponents NewConfig error: %v", err)
	}

	emailConfig, err := cfg.EmailConfig()
	if err != nil {
		return nil
	}

	return &mail.EmailInfo{
		Sender:    emailConfig.Sender,
		Password:  emailConfig.Password,
		Host:      emailConfig.Host,
		Port:      emailConfig.Port,
		Recipient: email,
		Subject:   mail.EmailSubject,
		Body:      mail.EmailBody,
		URL:       mail.EmailContentLink,
	}
}