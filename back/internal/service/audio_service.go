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

	audio.NeedConvert = audioFullPath == ""
	audio.HasAudio = audioFullPath != ""
	audio.AudioFullFilePath = &audioFullPath

	return s.audioRep.Update(ctx, audio)
}

func (s *AudioService) GetForConvert(ctx context.Context) ([]model.Audio, error) {
	return s.audioRep.GetForConvert(ctx)
}

func (s *AudioService) GetForTranscribe(ctx context.Context) ([]model.Audio, error) {
	return s.audioRep.GetForTranscribe(ctx)
}

//func (s *AudioService) SaveTranscribe(ctx context.Context, id string, transcribe string) error {
//	audio, err := s.GetByID(ctx, id)
//	if err != nil {
//		return err
//	}
//
//	audio.Transcribed = &transcribe
//	audio.HasTranscribed = true
//
//	return s.audioRep.Update(ctx, audio)
//
//}

func (s *AudioService) SaveTranscribeFullPath(ctx context.Context, id string, transcribeFullPath string) error {
	audio, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	audio.TranscribedFullFilePath = &transcribeFullPath
	audio.HasTranscribed = transcribeFullPath != ""

	return s.audioRep.Update(ctx, audio)
}

func (s *AudioService) GetTranscribeToSend(ctx context.Context) ([]model.Audio, error) {
	return s.audioRep.GetTranscribeToSend(ctx)
}

func (s *AudioService) SaveBotHubMessageID(ctx context.Context, id string, botHubMessageID string) error {
	audio, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	audio.BotHubMessageID = &botHubMessageID
	audio.IsMessageSend = botHubMessageID != ""

	return s.audioRep.Update(ctx, audio)
}

func (s *AudioService) GetMessagesForRead(ctx context.Context) ([]model.Audio, error) {
	return s.audioRep.GetMessagesForRead(ctx)
}

func (s *AudioService) SaveMessageTextFullPath(ctx context.Context, id string, messageTextFullFilePath string) error {
	audio, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	audio.MessageTextFullFilePath = &messageTextFullFilePath
	audio.IsMessageRead = messageTextFullFilePath != ""

	return s.audioRep.Update(ctx, audio)
}
