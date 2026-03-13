package service

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const audioUrl = "openai/v1/audio/transcriptions"
const chatUrl = "chat"
const messageUrl = "message"

//const botHubToken = "<ваш токен доступа>"

type BotHubService struct {
	url         string
	token       string
	fileService *FileService
}

func NewBotHubService(url string, token string, fileService *FileService) *BotHubService {
	return &BotHubService{
		url:         url,
		token:       token,
		fileService: fileService,
	}
}

func (s *BotHubService) NewChat(filePath string) error {
	//model.NewChatModel{
	//	ModelId:   "gpt",
	//	Name:      "Информация о сервисах с Архиектурного совета",
	//	Highlight: "#1C64F2",
	//	Platform:  "WEB",
	//}

	return nil
}

func (s *BotHubService) ProcessFileTranscribe() {
	files := s.fileService.GetAudioFiles()

	for _, file := range files {
		err := s.SendFileToTranscribe(file)
		if err != nil {
			log.Printf("convert error: %v", err)
			continue
		}

		err = s.fileService.MoveToEndedAudio(file)
		if err != nil {
			log.Printf("move error: %v", err)
		}
	}
}

func (s *BotHubService) SendFileToTranscribe(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(name))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.WriteField("model", "whisper-1")
	if err != nil {
		return err
	}

	err = writer.WriteField("response_format", "text")
	if err != nil {
		return err
	}

	//err = writer.WriteField("response_format", "verbose_json")
	//if err != nil {
	//	return err
	//}
	//
	//err = writer.WriteField("timestamp_granularities[]", "word")
	//if err != nil {
	//	return err
	//}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.url+audioUrl, &body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed: %s", resp.Status)
	}

	// формируем имя txt файла
	base := filepath.Base(name)
	nameWithoutExt := strings.TrimSuffix(base, filepath.Ext(base))
	txtPath := filepath.Join(filepath.Dir(name), nameWithoutExt+".txt")

	// сохраняем ответ
	err = os.WriteFile(txtPath, respBody, 0644)
	if err != nil {
		return err
	}

	_ = s.fileService.MoveToTranscribe(txtPath)

	fmt.Printf("transcription saved: %s\n", txtPath)

	return nil
}
