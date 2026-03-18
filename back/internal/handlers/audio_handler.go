package handlers

import (
	"analize_audio/internal/model"
	"analize_audio/internal/service"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AudioHandler struct {
	audioService *service.AudioService
	fileService  *service.FileService
}

func NewAudioHandler(audioService *service.AudioService, fileService *service.FileService) AudioHandler {
	return AudioHandler{
		audioService: audioService,
		fileService:  fileService,
	}
}

func (h AudioHandler) Get(w http.ResponseWriter, r *http.Request) {
	audios, err := h.audioService.Get(r.Context())
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
	}

	WriteJSON(w, http.StatusOK, audios)
}

func (h *AudioHandler) Add(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<30)

	reader, err := r.MultipartReader()
	if err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid multipart data")
		return
	}

	var (
		id           = uuid.New().String()
		fullFilePath string
		fileName     string
		thems        []string
	)

	dirSave := "converts"

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			if fullFilePath != "" {
				_ = h.fileService.RemoveByFullPath(r.Context(), fullFilePath)
			}

			WriteError(w, http.StatusInternalServerError, "Error reading multipart")

			return
		}

		switch part.FormName() {

		case "video":
			contentType := part.Header.Get("Content-Type")

			if !strings.HasPrefix(contentType, "video/") {
				_ = part.Close()
				continue
			}

			fileName = part.FileName()
			fullFilePath, _ = h.fileService.SaveFile(part, id, dirSave)

			_ = part.Close()

		case "thems":
			data, _ := io.ReadAll(part)
			them := strings.TrimSpace(string(data))
			if them != "" {
				thems = append(thems, them)
			}
			_ = part.Close()
		}

	}

	if fullFilePath == "" {
		WriteError(w, http.StatusBadRequest, "Отсутствует файл")
	}

	audio := model.Audio{
		ID:              id,
		SrcFullFilePath: fullFilePath,
		SrcFileName:     fileName,
		NeedConvert:     true,
		HasAudio:        false,
		HasTranscribed:  false,
		CreatedAt:       time.Now(),
	}

	if err := h.audioService.Create(r.Context(), audio, thems); err != nil {
		_ = h.fileService.RemoveByFullPath(r.Context(), fullFilePath)
		WriteError(w, http.StatusInternalServerError, "Не удалось сохранить фильм")
		return
	}

	WriteJSON(w, http.StatusOK, nil)
}

func (h *AudioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	audio, err := h.audioService.GetByID(r.Context(), id)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err.Error())
	}

	var analysisText string
	var transcribedText string

	if audio.IsMessageRead {
		if audio.MessageTextFullFilePath != nil {
			fileText, err := os.ReadFile(*audio.MessageTextFullFilePath)
			if err == nil {
				analysisText = string(fileText)
			}
		}
	}

	if audio.HasTranscribed {
		if audio.TranscribedFullFilePath != nil {
			fileText, err := os.ReadFile(*audio.TranscribedFullFilePath)
			if err == nil {
				transcribedText = string(fileText)
			}
		}
	}

	payload := map[string]interface{}{
		"id":               audio.ID,
		"src_file_name":    audio.SrcFileName,
		"analysis_text":    analysisText,
		"transcribed_text": transcribedText,
	}

	WriteJSON(w, http.StatusOK, payload)
}
