package logic

import (
	"context"
	"errors"

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
	if err != nil {
		return nil, errors.New("invalid access token")
	}

	userID, ok := claims["userID"].(int64)
	if !ok {
		return nil, errors.New("invalid user ID")
	}

	user, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(userID))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("user not found")
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
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
