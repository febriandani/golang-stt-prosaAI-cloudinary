package searchimage

import (
	"context"
	"database/sql"

	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/image"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/infra"
	"github.com/sirupsen/logrus"
)

type ImageConfig struct {
	db  *infra.DatabaseList
	log *logrus.Logger
}

func newDatabaseImage(db *infra.DatabaseList, logger *logrus.Logger) ImageConfig {
	return ImageConfig{
		db:  db,
		log: logger,
	}
}

type Image interface {
	GetImages(ctx context.Context, keyword string, offset, limit int) ([]image.ResponseSearchImages, error)
}

func (ic ImageConfig) GetImages(ctx context.Context, keyword string, offset, limit int) ([]image.ResponseSearchImages, error) {
	var result []image.ResponseSearchImages

	queryStatement := `select 
	si.id ,
	si.category_id ,
	si.keyword ,
	si.url 
	from search_images si where si.keyword ILIKE '%' || ? || '%' OFFSET ((? - 1) * ?) ROWS FETCH NEXT ? ROWS ONLY`

	query, args, err := ic.db.Backend.Read.In(queryStatement, keyword, offset, limit, limit)
	if err != nil {
		return nil, err
	}

	query = ic.db.Backend.Read.Rebind(query)
	err = ic.db.Backend.Read.Select(&result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}
