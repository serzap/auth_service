package logic

import (
	"context"
	"database/sql"

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
	user, err := l.svcCtx.UsersModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil {
		if err == model.ErrNotFound {
			logx.Errorf("User with email %s not found", in.Email)
			return &api.VerifyEmailResponse{Success: false}, ErrUserNotFound
		}
		return &api.VerifyEmailResponse{Success: false}, err
	}

	if user.VerificationCode.String != in.VerificationCode || !user.VerificationCode.Valid {
		logx.Errorf("Invalid verification code for user %s", in.Email)
		return &api.VerifyEmailResponse{Success: false}, ErrInvalidVerificationCode
	}

	user.Verified = true
	user.VerificationCode = sql.NullString{String: "", Valid: false}
	err = l.svcCtx.UsersModel.Update(l.ctx, user)
	if err != nil {
		return &api.VerifyEmailResponse{Success: false}, err
	}

	logx.Infof("User %s successfully verified", in.Email)
	return &api.VerifyEmailResponse{Success: true}, nil
}
