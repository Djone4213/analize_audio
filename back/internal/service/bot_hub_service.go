package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const audioUrl = "openai/v1/audio/transcriptions"
const messageIDUrl = "message/"
const messageSendUrl = "message/send"

type BotHubService struct {
	url                   string
	token                 string
	audioService          *AudioService
	outputTranscribeDir   string
	outputContentDir      string
	messagePromt          string
	messageAnswerTemplate string
	//fileService *FileService
}

func NewBotHubService(url string, token string, audioService *AudioService, outputTranscribeDir string, outputContentDir string) *BotHubService {
	promt, err := os.ReadFile("D:\\Work\\Projects\\analize_audio\\back\\promt.txt")

	var messagePromt string

	if err == nil {
		messagePromt = strings.TrimSpace(string(promt))
	}

	answerTemplate, err := os.ReadFile("D:\\Work\\Projects\\analize_audio\\back\\template_query.txt")

	var messageAnswerTemplate string

	if err == nil {
		messageAnswerTemplate = strings.TrimSpace(string(answerTemplate))
	}

	return &BotHubService{
		url:                   url,
		token:                 token,
		audioService:          audioService,
		outputTranscribeDir:   outputTranscribeDir,
		outputContentDir:      outputContentDir,
		messagePromt:          messagePromt,
		messageAnswerTemplate: messageAnswerTemplate,
	}
}

func (s *BotHubService) NewChat(filePath string) error {
	//chatModel := model.ChatModel{
	//	ModelId:   "gpt",
	//	Name:      "Информация о сервисах с Архитектурного совета",
	//	Highlight: "#1C64F2",
	//	Platform:  "WEB",
	//}

	return nil
}

func (s *BotHubService) ProcessSendMessages() {
	audios, err := s.audioService.GetTranscribeToSend(context.Background())

	if err != nil {
		log.Printf("ProcessSendMessages error: %v", err)
		return
	}

	log.Printf("⏱ Start send messages files count file:%d", len(audios))

	for _, audio := range audios {
		if audio.TranscribedFullFilePath == nil {
			continue
		}

		messageID, err := s.SendMessage(*audio.TranscribedFullFilePath)
		if err != nil {
			log.Printf("Send Message error: %v", err)
		}

		_ = s.audioService.SaveBotHubMessageID(context.Background(), audio.ID, messageID)
	}
}

func (s *BotHubService) ProcessFileTranscribe() {
	audios, err := s.audioService.GetForTranscribe(context.Background())

	if err != nil {
		log.Printf("ProcessFileTranscribe error: %v", err)
		return
	}

	log.Printf("⏱ Start transcribing files count file:%d", len(audios))

	for _, audio := range audios {
		if audio.AudioFullFilePath == nil {
			continue
		}
		transcribeFullPath, err := s.SendFileToTranscribe(*audio.AudioFullFilePath)
		if err != nil {
			log.Printf("transcribe error: %v", err)
			continue
		}

		_ = s.audioService.SaveTranscribeFullPath(context.Background(), audio.ID, transcribeFullPath)
	}
}

func (s *BotHubService) ProcessReadMessages() {
	audios, err := s.audioService.GetMessagesForRead(context.Background())

	if err != nil {
		log.Printf("ProcessReadMessages error: %v", err)
		return
	}

	log.Printf("⏱ Start read messages files count file:%d", len(audios))

	for _, audio := range audios {
		if audio.BotHubMessageID == nil {
			continue
		}

		messageTextFullPath, err := s.ReadMessage(*audio.BotHubMessageID, audio.ID)
		if err != nil {
			log.Printf("transcribe error: %v", err)
			continue
		}

		_ = s.audioService.SaveMessageTextFullPath(context.Background(), audio.ID, messageTextFullPath)
	}
}

func (s *BotHubService) ReadMessage(messageID string, filename string) (string, error) {
	if messageID == "" {
		return "", fmt.Errorf("messageID is empty")
	}

	req, err := http.NewRequest("GET", s.url+messageIDUrl+messageID, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed: %s, body: %s", resp.Status, string(respBody))
	}

	// структура ответа
	var result struct {
		Content string `json:"content"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w, body: %s", err, string(respBody))
	}

	if result.Content == "" {
		return "", fmt.Errorf("empty content in response: %s", string(respBody))
	}

	// формируем путь к файлу (по messageID)
	filename = fmt.Sprintf("%s.txt", filename)
	outputPath := filepath.Join(s.outputContentDir, filename)

	// сохраняем content
	err = os.WriteFile(outputPath, []byte(result.Content), 0644)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func (s *BotHubService) SendMessage(messageFilePath string) (string, error) {
	if s.messagePromt == "" {
		return "", fmt.Errorf("messagePromt file is empty")
	}

	if s.messageAnswerTemplate == "" {
		return "", fmt.Errorf("messageAnswerTemplate file is empty")
	}

	data, err := os.ReadFile(messageFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read message file: %w", err)
	}

	messageData := strings.TrimSpace(string(data))
	if messageData == "" {
		return "", fmt.Errorf("message file is empty")
	}

	messageText := fmt.Sprintf("%s\n%s\n%s", s.messagePromt, messageData, s.messageAnswerTemplate)

	message := map[string]interface{}{
		"chatId":   "076d6306-eac5-4a44-bad3-83022deac37f",
		"message":  messageText,
		"model_id": "gpt",
		//"sourceUserMessageId": "",
		//"tgBotMessageId":      "",
		"platform": "MAIN",
	}

	// сериализация в JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	// формируем запрос
	req, err := http.NewRequest("POST", s.url+messageSendUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("request failed: %s, body: %s", resp.Status, string(respBody))
	}

	//// путь для сохранения ответа
	//outputPath := filepath.Join(s.outputTranscribeDir, "message_response.json")
	//
	//err = os.WriteFile(outputPath, respBody, 0644)
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Printf("response saved: %s\n", outputPath)

	// структура для парсинга ответа
	var result struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w, body: %s", err, string(respBody))
	}

	if result.ID == "" {
		return "", fmt.Errorf("empty id in response: %s", string(respBody))
	}

	return result.ID, nil
}

func (s *BotHubService) SendFileToTranscribe(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := file.Close(); err != nil {
			// можно залогировать
		}
	}()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}

	err = writer.WriteField("model", "whisper-1")
	if err != nil {
		return "", err
	}

	err = writer.WriteField("response_format", "text")
	if err != nil {
		return "", err
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
		return "", err
	}

	req, err := http.NewRequest("POST", s.url+audioUrl, &body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed: %s", resp.Status)
	}

	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	outputPath := filepath.Join(s.outputTranscribeDir, name+".txt")

	// сохраняем ответ
	err = os.WriteFile(outputPath, respBody, 0644)
	if err != nil {
		return "", err
	}

	fmt.Printf("transcription saved: %s\n", outputPath)

	return outputPath, nil
}
