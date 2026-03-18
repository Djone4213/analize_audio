package repositories

import (
	"analize_audio/internal/model"

	"gorm.io/gorm"
)

func Migrate(DB *gorm.DB) {
	if err := DB.AutoMigrate(
		&model.Audio{},
	); err != nil {
		panic(err)
	}
}
