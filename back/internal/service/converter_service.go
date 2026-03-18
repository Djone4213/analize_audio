package service

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type ConverterService struct {
	fileService *FileService
}

func NewConverterService(fileService *FileService) *ConverterService {
	return &ConverterService{
		fileService: fileService,
	}
}

func (s *ConverterService) ProcessConvertFiles() {
	files := s.fileService.GetMKVFiles()

	for _, file := range files {
		err := s.ConvertToAudio(file)
		if err != nil {
			log.Printf("convert error: %v", err)
			continue
		}

		err = s.fileService.MoveToConverted(file)
		if err != nil {
			log.Printf("move error: %v", err)
		}
	}
}

func (s *ConverterService) ConvertToAudio(filename string) error {
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	outputPath := filepath.Join(s.fileService.GetDir(), audioDir, name+".mp3")

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
		return err
	}

	log.Printf("converted: %s -> %s", filename, outputPath)

	return nil
}
