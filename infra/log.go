package infra

import (
	"os"

	constants "github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/constants/general"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/utils"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func NewLogger(conf *general.AppService) *logrus.Logger {
	if logger == nil {
		path := "log/"

		isExist, err := utils.DirExists(path)
		if err != nil {
			panic(err)
		}

		if !isExist {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		writer, err := rotatelogs.New(
			path+conf.App.Name+"-"+"%Y%m%d.log",
			rotatelogs.WithMaxAge(-1),
			rotatelogs.WithRotationCount(constants.MaxRotationFile),
			rotatelogs.WithRotationTime(constants.LogRotationTime),
		)
		if err != nil {
			panic(err)
		}

		logger = logrus.New()

		// TODO: Active this code if later it's needed to limit the log
		// // Set Log level that need to show or stored
		// if conf.App.Environtment == constants.EnvProd {
		// 	logger.SetLevel(logrus.WarnLevel)
		// } else {
		// 	logger.SetLevel(logrus.DebugLevel)
		// }

		// Set Hook with writer & formatter for log file
		logger.Hooks.Add(lfshook.NewHook(
			writer,
			&logrus.TextFormatter{
				DisableColors:   false,
				FullTimestamp:   true,
				TimestampFormat: constants.FullTimeFormat,
			},
		))

		// Set formatter for os.Stdout
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: constants.FullTimeFormat,
		})

		return logger
	}

	return logger
}

func TestNewLogger() *logrus.Logger {
	logger := logrus.New()
	return logger
}
