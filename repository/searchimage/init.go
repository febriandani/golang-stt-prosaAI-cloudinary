package searchimage

import (
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/infra"
	"github.com/sirupsen/logrus"
)

type DatabaseImage struct {
	Image Image
}

func NewDatabaseImage(db *infra.DatabaseList, logger *logrus.Logger) DatabaseImage {
	return DatabaseImage{
		Image: newDatabaseImage(db, logger),
	}
}
