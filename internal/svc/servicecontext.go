package svc

import (
	"os"

	"github.com/serzap/auth_service/internal/config"
	"github.com/serzap/auth_service/model"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config     config.Config
	UsersModel model.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	logx.Info("try open database")
	conn := sqlx.NewMysql(c.DataSource)
	sqlContent, err := os.ReadFile(c.SqlFile)
	if err != nil {
		logx.Info("Failed read content")
	}
	res, err := conn.Exec(string(sqlContent))
	if err != nil {
		logx.Info("Migrations is not applied")
	}
	logx.Info(res)
	usersModel := model.NewUsersModel(conn)
	return &ServiceContext{
		Config:     c,
		UsersModel: usersModel,
	}
}
