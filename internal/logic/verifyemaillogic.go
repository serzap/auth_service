package logic

import (
	"context"
	"database/sql"
	"errors"

	"github.com/serzap/auth_service/api"
	"github.com/serzap/auth_service/internal/svc"
	"github.com/serzap/auth_service/model"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
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
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if user.VerificationCode.String != in.VerificationCode || !user.VerificationCode.Valid {
		logx.Errorf("Invalid verification code for user %s", in.Email)
		return nil, ErrInvalidVerificationCode
	}

	user.Verified = true
	user.VerificationCode = sql.NullString{String: "", Valid: false} // Очистить verification code
	err = l.svcCtx.UsersModel.Update(l.ctx, user)
	if err != nil {
		return nil, err
	}

	logx.Infof("User %s successfully verified", in.Email)
	return &api.VerifyEmailResponse{}, nil
}
