package logic

import (
	"context"

	"github.com/serzap/auth_service/api"
	"github.com/serzap/auth_service/internal/svc"

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
	// todo: add your logic here and delete this line

	return &api.GetUserInfoResponse{}, nil
}
