package repository

import (
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/infra"
	ri "github.com/febriandani/golang-stt-prosaAI-cloudinary/repository/searchimage"
	ru "github.com/febriandani/golang-stt-prosaAI-cloudinary/repository/user"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	DatabaseUser  ru.DatabaseUser
	DatabaseImage ri.DatabaseImage
}

func NewRepo(database *infra.DatabaseList, logger *logrus.Logger) Repo {
	return Repo{
		DatabaseUser:  ru.NewDatabaseUser(database, logger),
		DatabaseImage: ri.NewDatabaseImage(database, logger),
	}
}
