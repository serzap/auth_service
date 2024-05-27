package logic

import (
	"context"

	"github.com/serzap/auth_service/api"
	"github.com/serzap/auth_service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyEmailLogic {
	return &VerifyEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyEmailLogic) VerifyEmail(in *api.VerifyEmailRequest) (*api.VerifyEmailResponse, error) {
	// todo: add your logic here and delete this line

	return &api.VerifyEmailResponse{}, nil
}
