package user

import (
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/service"
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
