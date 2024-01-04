package service

import (
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/infra"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/repository"
	su "github.com/febriandani/golang-stt-prosaAI-cloudinary/service/user"
	"github.com/go-redis/redis/v8"
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
