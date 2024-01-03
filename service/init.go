package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/pharmaniaga/auth-user/domain/model/general"
	"github.com/pharmaniaga/auth-user/infra"
	"github.com/pharmaniaga/auth-user/repository"
	su "github.com/pharmaniaga/auth-user/service/user"
	"github.com/sirupsen/logrus"
)

type Service struct {
	User su.ServiceUser
}

func NewService(repo repository.Repo, conf general.AppService, dbList *infra.DatabaseList, logger *logrus.Logger, redis *redis.Client) Service {
	return Service{
		User: su.NewServiceUser(repo.DatabaseUser, logger, conf, dbList, redis),
	}
}
