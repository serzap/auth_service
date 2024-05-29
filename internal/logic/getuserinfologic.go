package logic

import (
	"context"
	"fmt"

	"github.com/serzap/auth_service/api"
	"github.com/serzap/auth_service/internal/svc"
	"github.com/serzap/auth_service/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *api.GetUserInfoRequest) (*api.GetUserInfoResponse, error) {
	claims, err := l.parseToken(in.Token)
	logx.Info(claims)
	if err != nil {
		logx.Error(err)
		return nil, ErrInvalidToken
	}

	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		logx.Error(err)
		return nil, ErrInvalidUserID
	}
	userID := int64(userIDFloat)

	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(userID))
	if err != nil {
		logx.Error(err)
		if err == model.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	response := &api.GetUserInfoResponse{
		UserId:    userID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	return response, nil
}

func (l *GetUserInfoLogic) parseToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := []byte(l.svcCtx.Config.SecretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
