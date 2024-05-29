package logic

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"regexp"

	"github.com/serzap/auth_service/api"
	"github.com/serzap/auth_service/internal/svc"
	"github.com/serzap/auth_service/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *api.RegisterRequest) (*api.RegisterResponse, error) {
	logx.Info("start registering")
	if !isValidEmail(in.Email) {
		return nil, errors.New("invalid email format")
	}
	logx.Info("succesfull email validation")

	if exists, err := l.isUsernameOrEmailExists(in.Username, in.Email); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.New("username or email already exists")
	}

	user := &model.Users{
		Email:     in.Email,
		Username:  in.Username,
		FirstName: in.FirstName,
		LastName:  in.LastName,
	}
	result, err := l.svcCtx.UsersModel.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	go l.sendVerificationEmail(user.Email)
	return &api.RegisterResponse{UserId: userID}, nil
}

func isValidEmail(email string) bool {
	// Define a regular expression for validating an email address.
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Функция для проверки существования username или email
func (l *RegisterLogic) isUsernameOrEmailExists(username, email string) (bool, error) {
	user, err := l.svcCtx.UsersModel.FindOneByEmail(l.ctx, email)
	if err != nil && err != model.ErrNotFound {
		return false, err
	}
	if user != nil {
		return true, nil
	}

	user, err = l.svcCtx.UsersModel.FindOneByUsername(l.ctx, username)
	if err != nil && err != model.ErrNotFound {
		return false, err
	}
	if user != nil {
		return true, nil
	}

	return false, nil
}

const (
	smtpServer = "smtp.example.com"         // Замените на адрес вашего SMTP-сервера
	smtpPort   = "587"                      // Порт вашего SMTP-сервера
	smtpUser   = "auth_service@example.com" // Ваш SMTP логин
	smtpPass   = "1111"                     // Ваш SMTP пароль
)

func (l *RegisterLogic) sendVerificationEmail(email string) {
	verificationCode := generateVerificationCode()
	if verificationCode == "" {
		logx.Errorf("Failed to generate verification code for email: %s", email)
		return
	}
	// Set up authentication information.
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpServer)

	// Message body.
	subject := "Subject: Email Verification Code\n"
	body := fmt.Sprintf("Your verification code is: %s", verificationCode)
	msg := []byte(subject + "\n" + body)

	// Sending email.
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpUser, []string{email}, msg)
	if err != nil {
		log.Printf("Failed to send verification email to %s: %v", email, err)
		return
	}
	log.Printf("Verification email sent to %s", email)
	logx.Infof("Sending verification email to %s with code %s", email, verificationCode)
}

const verificationCodeLength = 6
const digits = "0123456789"

func generateVerificationCode() string {
	code := make([]byte, verificationCodeLength)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			logx.Errorf("Error generating verification code: %v", err)
			return ""
		}
		code[i] = digits[num.Int64()]
	}
	return string(code)
}
