package api

import (
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/handler/api/authorization"
	hu "github.com/febriandani/golang-stt-prosaAI-cloudinary/handler/api/user"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/service"
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
