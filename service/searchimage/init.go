package searchimage

import (
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/infra"
	ri "github.com/febriandani/golang-stt-prosaAI-cloudinary/repository/searchimage"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type ServiceImage struct {
	Image Image
}

func NewServiceImage(database ri.DatabaseImage, logger *logrus.Logger, conf general.AppService, dbList *infra.DatabaseList, redis *redis.Client) ServiceImage {
	return ServiceImage{
		Image: newImageService(database, logger, dbList, conf, redis),
	}
}
