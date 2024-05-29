package logic

import (
	"context"
	"errors"

	"github.com/serzap/auth_service/api"
	"github.com/serzap/auth_service/internal/svc"
	"github.com/serzap/auth_service/model"

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
	user, err := l.svcCtx.UsersModel.FindByVerificationCode(l.ctx, in.VerificationCode)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errors.New("verification code not found")
		}
		return nil, err
	}

	valid, err := l.svcCtx.UsersModel.IsVerificationCodeValid(l.ctx, in.VerificationCode)

	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, err
	}

	user.Verified = true
	err = l.svcCtx.UsersModel.Update(l.ctx, user)
	if err != nil {
		return nil, err
	}

	return &api.VerifyEmailResponse{}, nil
}
