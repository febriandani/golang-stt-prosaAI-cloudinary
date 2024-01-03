package api

import (
	"github.com/pharmaniaga/auth-user/domain/model/general"
	"github.com/pharmaniaga/auth-user/handler/api/authorization"
	hu "github.com/pharmaniaga/auth-user/handler/api/user"
	"github.com/pharmaniaga/auth-user/service"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Token  authorization.TokenHandler
	Public authorization.PublicHandler
	User   hu.HandlerUser
}

func NewHandler(sv service.Service, conf general.AppService, logger *logrus.Logger) Handler {
	return Handler{
		Token:  authorization.NewTokenHandler(conf, logger),
		Public: authorization.NewPublicHandler(conf, logger),
		User:   hu.NewHandlerUser(sv, conf, logger),
	}
}