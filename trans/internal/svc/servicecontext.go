package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"tcc/model"
	"tcc/trans/internal/config"
)

type ServiceContext struct {
	Config           config.Config
	UserAccountModel model.UserAccountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysql := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:           c,
		UserAccountModel: model.NewUserAccountModel(mysql, c.CacheRedis),
	}
}
