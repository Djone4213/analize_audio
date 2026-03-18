package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileService struct {
	dir string
}

func NewFileService(dir string) *FileService {
	return &FileService{dir: dir}
}

func (s *FileService) RemoveByFullPath(ctx context.Context, fullFilePath string) error {
	return os.Remove(fullFilePath)
}

// SaveFile сохраняет файл и возвращает имя и полный путь
func (s *FileService) SaveFile(part *multipart.Part, fileName string, dirSave string) (string, error) {
	fileName = fileName + filepath.Ext(part.FileName())

	filePath := filepath.Join(s.dir, dirSave, fileName)

	if err := s.CreateFullPath(filepath.Dir(filePath)); err != nil {
		return "", fmt.Errorf("не удалось создать каталог: %w", err)
	}

	if _, err := os.Stat(filePath); err == nil {
		_ = os.Remove(filePath)
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer func() {
		_ = dst.Close()
	}()

	//time.Sleep(100 * time.Millisecond) // Для Dev
	if _, err := io.Copy(dst, part); err != nil {
		_ = os.Remove(filePath)
		return "", fmt.Errorf("не удалось записать файл: %w", err)
	}

	return filePath, nil
}

func (s *FileService) CreateFullPath(filePath string) error {
	return os.MkdirAll(filePath, 0755)
}

//func (s *FileService) GetDir() string {
//	return s.dir
//}
//
//func (s *FileService) GetMKVFiles() []string {
//	filesPath := make([]string, 0)
//
//	fullPath := filepath.Join(s.dir, convertDir)
//
//	files, err := os.ReadDir(fullPath)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, file := range files {
//		if !file.IsDir() {
//			filesPath = append(filesPath, filepath.Join(fullPath, file.Name()))
//		}
//	}
//
//	return filesPath
//}
//
//func (s *FileService) MoveToConverted(filename string) error {
//	convertedDir := filepath.Join(s.dir, converted)
//
//	return s.moveFile(filename, convertedDir)
//}
//
//func (s *FileService) MoveToEndedAudio(filename string) error {
//	movedDir := filepath.Join(s.dir, audioEndDir)
//
//	return s.moveFile(filename, movedDir)
//}
//
//func (s *FileService) MoveToTranscribe(filename string) error {
//	transcribeDir := filepath.Join(s.dir, transcribe)
//
//	return s.moveFile(filename, transcribeDir)
//}
//
//func (s *FileService) moveFile(filename string, moveDir string) error {
//	// создаём папку если нет
//	err := os.MkdirAll(moveDir, os.ModePerm)
//	if err != nil {
//		return err
//	}
//
//	dst := filepath.Join(moveDir, filepath.Base(filename))
//
//	return os.Rename(filename, dst)
//}
//
//func (s *FileService) GetAudioFiles() []string {
//	filesPath := make([]string, 0)
//
//	fullPath := filepath.Join(s.dir, audioDir)
//
//	files, err := os.ReadDir(fullPath)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, file := range files {
//		if !file.IsDir() {
//			filesPath = append(filesPath, filepath.Join(fullPath, file.Name()))
//		}
//	}
//
//	return filesPath
//}
