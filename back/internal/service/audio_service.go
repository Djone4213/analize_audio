package service

import (
	"analize_audio/internal/model"
	"analize_audio/internal/repositories"
	"context"
)

type AudioService struct {
	audioRep *repositories.AudioRepository
}

func NewAudioService(audioRep *repositories.AudioRepository) *AudioService {
	return &AudioService{
		audioRep: audioRep,
	}
}

func (s *AudioService) Create(ctx context.Context, audio model.Audio) error {
	return s.audioRep.Create(ctx, audio)
}

func (s *AudioService) Get(ctx context.Context) ([]model.Audio, error) {
	return s.audioRep.Get(ctx)
}

func (s *AudioService) GetByID(ctx context.Context, id string) (model.Audio, error) {
	return s.audioRep.GetByID(ctx, id)
}

func (s *AudioService) SaveAudioFullPath(ctx context.Context, id string, audioFullPath string) error {
	audio, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	audio.NeedConvert = audioFullPath != ""
	audio.HasAudio = audioFullPath != ""
	audio.AudioFullFilePath = &audioFullPath

	return s.audioRep.Update(ctx, audio)
}
