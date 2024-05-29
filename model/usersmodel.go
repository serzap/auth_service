package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		withSession(session sqlx.Session) UsersModel
		FindByVerificationCode(ctx context.Context, code string) (*Users, error)
		IsVerificationCodeValid(ctx context.Context, code string) (bool, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn),
	}
}

func (m *customUsersModel) withSession(session sqlx.Session) UsersModel {
	return NewUsersModel(sqlx.NewSqlConnFromSession(session))
}

func (m *customUsersModel) FindByVerificationCode(ctx context.Context, code string) (*Users, error) {
	var user Users
	query := fmt.Sprintf("select %s from %s where `verification_code` = ? limit 1", usersRows, m.table)
	err := m.conn.QueryRowCtx(ctx, query, code)
	switch {
	case err == nil:
		return &user, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customUsersModel) IsVerificationCodeValid(ctx context.Context, code string) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where `verification_code` = ? and `verified` = ?", m.table)
	var count int
	err := m.conn.QueryRowCtx(ctx, query, code, true)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
