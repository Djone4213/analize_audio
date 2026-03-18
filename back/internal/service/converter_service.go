package service

import (
	"context"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type ConverterService struct {
	audioService   *AudioService
	outputAudioDir string
}

func NewConverterService(audioService *AudioService, outputAudioDir string) *ConverterService {
	return &ConverterService{
		audioService:   audioService,
		outputAudioDir: outputAudioDir,
	}
}

func (s *ConverterService) ProcessConvertFiles() {
	audios, err := s.audioService.GetForConvert(context.Background())

	if err != nil {
		log.Printf("convert error: %v", err)
		return
	}

	log.Printf("⏱ Start converts files count file:%d", len(audios))
	
	for _, audio := range audios {
		outputFilePath, err := s.ConvertToAudio(audio.SrcFullFilePath)
		if err != nil {
			log.Printf("convert error: %v", err)
			continue
		}

		_ = s.audioService.SaveAudioFullPath(context.Background(), audio.ID, outputFilePath)
	}

}

func (s *ConverterService) ConvertToAudio(filename string) (string, error) {
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	outputPath := filepath.Join(s.outputAudioDir, name+".mp3")

	cmd := exec.Command(
		"ffmpeg",
		"-i", filename,
		"-vn",
		"-acodec", "libmp3lame",
		//"-ab", "64k",
		"-ab", "32k",
		outputPath,
	)

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	log.Printf("converted: %s -> %s", filename, outputPath)

	return outputPath, nil
}
