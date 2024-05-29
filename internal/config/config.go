package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	DataSource string
	SqlFile    string
	SmtpServer string
	SmtpPort   string
	SmtpUser   string
	SmtpPass   string
}
