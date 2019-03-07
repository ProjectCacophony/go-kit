package config

import (
	"github.com/jinzhu/gorm"
)

type Item struct {
	gorm.Model
	Namespace string
	Key       string
	Value     []byte
}

func (*Item) TableName() string {
	return "config"
}
