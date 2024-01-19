package searchimage

import (
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/service"
	"github.com/sirupsen/logrus"
)

type HandlerImage struct {
	Image ImageHandler
}

func NewHandlerImage(sv service.Service, conf general.AppService, logger *logrus.Logger) HandlerImage {
	return HandlerImage{
		Image: newImageHandler(sv.Image, conf, logger),
	}
}
