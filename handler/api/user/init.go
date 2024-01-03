package user

import (
	"github.com/pharmaniaga/auth-user/domain/model/general"
	"github.com/pharmaniaga/auth-user/service"
	"github.com/sirupsen/logrus"
)

type HandlerUser struct {
	User UserHandler
}

func NewHandlerUser(sv service.Service, conf general.AppService, logger *logrus.Logger) HandlerUser {
	return HandlerUser{
		User: newUserHandler(sv.User, conf, logger),
	}
}
