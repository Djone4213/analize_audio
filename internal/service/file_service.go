package service

import (
	"log"
	"os"
	"path/filepath"
)

type FileService struct {
	Dir string
}

func NewFileService(dir string) *FileService {
	return &FileService{Dir: dir}
}

func (s *FileService) GetMKVFiles() []string {
	filesPath := make([]string, 0)

	fullPath := filepath.Join(s.Dir, convertDir)

	files, err := os.ReadDir(fullPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filesPath = append(filesPath, filepath.Join(fullPath, file.Name()))
		}
	}

	return filesPath
}

func (s *FileService) MoveToConverted(filename string) error {
	convertedDir := filepath.Join(s.Dir, converted)

	return s.moveFile(filename, convertedDir)
}

func (s *FileService) MoveToEndedAudio(filename string) error {
	movedDir := filepath.Join(s.Dir, audioEndDir)

	return s.moveFile(filename, movedDir)
}

func (s *FileService) MoveToTranscribe(filename string) error {
	transcribeDir := filepath.Join(s.Dir, transcribe)

	return s.moveFile(filename, transcribeDir)
}

func (s *FileService) moveFile(filename string, moveDir string) error {
	// создаём папку если нет
	err := os.MkdirAll(moveDir, os.ModePerm)
	if err != nil {
		return err
	}

	dst := filepath.Join(moveDir, filepath.Base(filename))

	return os.Rename(filename, dst)
}

func (s *FileService) GetAudioFiles() []string {
	filesPath := make([]string, 0)

	fullPath := filepath.Join(s.Dir, audioDir)

	files, err := os.ReadDir(fullPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filesPath = append(filesPath, filepath.Join(fullPath, file.Name()))
		}
	}

	return filesPath
}
