package repositories

import (
	"analize_audio/internal/model"
	"context"

	"gorm.io/gorm"
)

type AudioRepository struct {
	db *gorm.DB
}

func NewAudioRepository(db *gorm.DB) *AudioRepository {
	return &AudioRepository{db: db}
}

func (r *AudioRepository) Create(ctx context.Context, audio model.Audio) error {
	return r.db.Create(&audio).Error
}

func (r *AudioRepository) Get(ctx context.Context) ([]model.Audio, error) {
	var audios []model.Audio
	err := r.db.Find(&audios).Error
	return audios, err
}

func (r *AudioRepository) GetByID(ctx context.Context, id string) (model.Audio, error) {
	var audio model.Audio
	err := r.db.Where("id = ?", id).First(&audio).Error
	return audio, err
}

func (r *AudioRepository) Update(ctx context.Context, audio model.Audio) error {
	return r.db.Save(&audio).Error
}

func (r *AudioRepository) GetForConvert(ctx context.Context) ([]model.Audio, error) {
	var audios []model.Audio
	err := r.db.
		Where("need_convert = true").
		Find(&audios).Error
	return audios, err
}

func (r *AudioRepository) GetForTranscribe(ctx context.Context) ([]model.Audio, error) {
	var audios []model.Audio
	err := r.db.
		Where("has_transcribed = false and has_audio = true").
		Find(&audios).Error
	return audios, err
}

func (r *AudioRepository) GetTranscribeToSend(ctx context.Context) ([]model.Audio, error) {
	var audios []model.Audio
	err := r.db.
		Where("has_transcribed = true and is_message_send = false").
		Find(&audios).Error
	return audios, err
}

func (r *AudioRepository) GetMessagesForRead(ctx context.Context) ([]model.Audio, error) {
	var audios []model.Audio
	err := r.db.
		Where("is_message_send = true and is_message_read = false").
		Find(&audios).Error
	return audios, err
}
