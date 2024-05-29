package logic

import (
	"context"
	"time"

	"github.com/serzap/auth_service/api"
	"github.com/serzap/auth_service/internal/svc"
	"github.com/serzap/auth_service/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *api.LoginRequest) (*api.LoginResponse, error) {
	user, err := l.svcCtx.UsersModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil {
		if err == model.ErrNotFound {
			logx.Errorf("Email not found %s", in.Email)
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(in.Password))
	if err != nil {
		logx.Errorf("hash and password comparison for %s failed", in.Email)
		return nil, ErrInvalidCredentials
	}
	token, err := generateJWT(user, l.svcCtx.Config.JWTExpiration, l.svcCtx.Config.SecretKey)
	if err != nil {
		logx.Errorf("JWT generation for %s failed with error: %w", in.Email, err)
		return nil, err
	}
	return &api.LoginResponse{
		Token: token,
	}, nil
}

func generateJWT(user *model.Users, expiration time.Duration, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"userID": user.Id,
		"email":  user.Email,
		"exp":    time.Now().Add(expiration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
