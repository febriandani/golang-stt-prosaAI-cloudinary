package user

import (
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/infra"
	ru "github.com/febriandani/golang-stt-prosaAI-cloudinary/repository/user"
	"github.com/go-redis/redis/v8"
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
