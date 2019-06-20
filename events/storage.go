package events

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"gocloud.dev/blob"
)

type Storage struct {
	logger *zap.Logger
	db     *gorm.DB
	bucket *blob.Bucket
}

func NewStorage(
	logger *zap.Logger,
	db *gorm.DB,
	bucket *blob.Bucket,
) *Storage {
	return &Storage{
		logger: logger,
		db:     db,
		bucket: bucket,
	}
}
