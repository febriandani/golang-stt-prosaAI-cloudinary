package repository

import (
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/infra"
	ru "github.com/febriandani/golang-stt-prosaAI-cloudinary/repository/user"
	"github.com/sirupsen/logrus"
)

type Repo struct {
	DatabaseUser ru.DatabaseUser
}

func NewRepo(database *infra.DatabaseList, logger *logrus.Logger) Repo {
	return Repo{
		DatabaseUser: ru.NewDatabaseUser(database, logger),
	}
}
