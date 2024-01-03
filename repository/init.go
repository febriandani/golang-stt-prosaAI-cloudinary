package repository

import (
	"github.com/pharmaniaga/auth-user/infra"
	ru "github.com/pharmaniaga/auth-user/repository/user"
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
