package db

import "gorm.io/gorm"

type GormManager struct {
	Orm *gorm.DB
}

var Manager *GormManager

func Setup(orm *gorm.DB) {
	Manager = &GormManager{
		Orm: orm,
	}
}
