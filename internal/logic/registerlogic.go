package logic

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"net/smtp"
	"regexp"

	"github.com/serzap/auth_service/api"
	"github.com/serzap/auth_service/internal/svc"
	"github.com/serzap/auth_service/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

const (
	codeLength    = 6
	digitsForCode = "0123456789"
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
		logx.Errorf("failed email validation")
		return nil, ErrInvalidEmail
	}
	logx.Info("succesfull email validation")

	exists, err := l.isUsernameOrEmailExists(in.Username, in.Email)
	if err != nil {
		logx.Errorf("failed email validation")
		return nil, err
	} else if exists {
		logx.Errorf("failed email validation")
		return nil, ErrInvalidEmail
	}
	logx.Info("succesfull credentials validation")

	passHash, err := hashPassword(in.Password)
	if err != nil {
		logx.Errorf("failed password hashing")
		return nil, err
	}
	logx.Info("succesfull password hashing")
	verificationCode, err := generateVerificationCode()
	if err != nil {
		logx.Errorf("failed to generate verification code")
		return nil, err
	}
	logx.Infof("Code for email %s generated", in.Email)

	user := &model.Users{
		Email:            in.Email,
		Username:         in.Username,
		FirstName:        in.FirstName,
		LastName:         in.LastName,
		PassHash:         passHash,
		VerificationCode: sql.NullString{String: verificationCode, Valid: true},
		Verified:         false,
	}
	result, err := l.svcCtx.UsersModel.Insert(l.ctx, user)
	if err != nil {
		return nil, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	go l.sendVerificationEmail(user.Email, user.VerificationCode.String)
	return &api.RegisterResponse{UserId: userID}, nil
}

func isValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

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

func generateVerificationCode() (string, error) {
	code := make([]byte, codeLength)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digitsForCode))))
		if err != nil {
			logx.Error(err)
			return "", err
		}
		code[i] = digitsForCode[num.Int64()]
	}
	return string(code), nil
}

func (l *RegisterLogic) sendVerificationEmail(email string, verificationCode string) {
	auth := smtp.PlainAuth("", l.svcCtx.Config.SmtpUser, l.svcCtx.Config.SmtpPass, l.svcCtx.Config.SmtpServer)

	subject := "Subject: Email Verification Code\n"
	body := fmt.Sprintf("Your verification code is: %s", verificationCode)
	msg := []byte(subject + "\n" + body)

	err := smtp.SendMail(l.svcCtx.Config.SmtpServer+":"+l.svcCtx.Config.SmtpPort, auth, l.svcCtx.Config.SmtpUser, []string{email}, msg)
	if err != nil {
		logx.Error(err)
		return
	}
	logx.Infof("Sending verification email to %s with code %s", email, verificationCode)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
