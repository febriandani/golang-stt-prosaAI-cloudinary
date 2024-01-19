package searchimage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	mu "github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/image"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/utils"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/infra"
	ri "github.com/febriandani/golang-stt-prosaAI-cloudinary/repository/searchimage"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type ImageService struct {
	db     ri.DatabaseImage
	log    *logrus.Logger
	conf   general.AppService
	dbConn *infra.DatabaseList
	redis  *redis.Client
}

func newImageService(database ri.DatabaseImage, logger *logrus.Logger, dbConn *infra.DatabaseList, conf general.AppService, redis *redis.Client) ImageService {
	return ImageService{
		db:     database,
		log:    logger,
		conf:   conf,
		dbConn: dbConn,
		redis:  redis,
	}
}

type Image interface {
	GetImages(ctx context.Context, filter mu.RequestSearchImage) (int64, interface{}, map[string]string, error)
}

func (is ImageService) GetImages(ctx context.Context, filter mu.RequestSearchImage) (int64, interface{}, map[string]string, error) {

	var responseImages []mu.ResponseSearchImages
	keys := make(map[int64]bool)

	tempText := utils.GenerateSubstrings(filter.Keyword)

	// Print the generated substrings
	for i, substring := range tempText {
		fmt.Printf("tempText[%d] = \"%s\"\n", i, substring)
		data, err := is.db.Image.GetImages(ctx, substring, filter.Offset, filter.Limit)
		if err != nil {
			log.Println("ERROR GET DATA Images", err.Error())
			return 0, nil, map[string]string{
				"en": "There was an error during get data users",
				"id": "Ada kesalahan saat menampilkan data staff",
			}, err
		}

		for _, val := range data {
			if _, valFlag := keys[val.ID]; !valFlag {
				keys[val.ID] = true
				responseImages = append(responseImages, val)

			}
		}

	}

	if len(responseImages) == 0 {
		return 0, nil, map[string]string{
			"en": "Data is empty",
			"id": "Data tidak ditemukan",
		}, errors.New("404Data")
	}

	return 0, responseImages, map[string]string{
		"en": "Successfully retrieved data users",
		"id": "Berhasil menampilkan data staff",
	}, nil
}
