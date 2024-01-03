package user

import (
	"github.com/pharmaniaga/auth-user/infra"
	"github.com/sirupsen/logrus"
)

type DatabaseUser struct {
	User User
}

func NewDatabaseUser(db *infra.DatabaseList, logger *logrus.Logger) DatabaseUser {
	return DatabaseUser{
		User: newDatabaseUser(db, logger),
	}
}
