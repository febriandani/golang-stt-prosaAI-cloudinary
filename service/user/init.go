package user

import (
	"github.com/go-redis/redis/v8"
	"github.com/pharmaniaga/auth-user/domain/model/general"
	"github.com/pharmaniaga/auth-user/infra"
	ru "github.com/pharmaniaga/auth-user/repository/user"
	"github.com/sirupsen/logrus"
)

type ServiceUser struct {
	User User
}

func NewServiceUser(database ru.DatabaseUser, logger *logrus.Logger, conf general.AppService, dbList *infra.DatabaseList, redis *redis.Client) ServiceUser {
	return ServiceUser{
		User: newUserService(database, logger, dbList, conf, redis),
	}
}
