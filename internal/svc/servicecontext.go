package svc

import (
	"github.com/serzap/auth_service/internal/config"
	"github.com/serzap/auth_service/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	UsersModel model.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	usersModel := model.NewUsersModel(conn)
	return &ServiceContext{
		Config:     c,
		UsersModel: usersModel,
	}
}
