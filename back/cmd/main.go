package main

import (
	"analize_audio/internal/config"
	"analize_audio/internal/handlers"
	"analize_audio/internal/middleware"
	"analize_audio/internal/repositories"
	"analize_audio/internal/service"
	"analize_audio/internal/worker"
	"analize_audio/pkg/db/gorm"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	cfg := config.Load()

	DB := gorm.InitDB(cfg.DB.Host, cfg.DB.Name, cfg.DB.User, cfg.DB.Password, cfg.DB.Port)
	repositories.Migrate(DB)

	aRep := repositories.NewAudioRepository(DB)

	fServ := service.NewFileService(cfg.App.Dir)
	aServ := service.NewAudioService(aRep)

	outputAudioDir := filepath.Join(cfg.App.Dir, "audio")
	outputTranscribeDir := filepath.Join(cfg.App.Dir, "transcribe")
	outputContentDir := filepath.Join(cfg.App.Dir, "content")

	if err := fServ.CreateFullPath(outputAudioDir); err != nil {
		log.Fatalf("Ошибка создания пути для сохранения аудио файлов: %v", err)
	}

	if err := fServ.CreateFullPath(outputTranscribeDir); err != nil {
		log.Fatalf("Ошибка создания пути для сохранения файлов транскрибации: %v", err)
	}

	if err := fServ.CreateFullPath(outputContentDir); err != nil {
		log.Fatalf("Ошибка создания пути для сохранения файлов с результатами ответов: %v", err)
	}

	cServ := service.NewConverterService(aServ, outputAudioDir)

	// 🔥 запускаем worker
	convWorker := worker.NewConverterWorker(cServ, time.Minute)
	convWorker.Start()

	bhServ := service.NewBotHubService(cfg.Bot.URL, cfg.Bot.Token, aServ, outputTranscribeDir, outputContentDir)

	//bhServ.ProcessSendMessages()

	// 🔥 worker для транскрибации
	transcribeWorker := worker.NewTranscribeWorker(bhServ, time.Minute)
	transcribeWorker.Start()

	// 🔥 worker для отправки сообщений
	sendMessagesWorker := worker.NewSendMessagesWorker(bhServ, time.Minute)
	sendMessagesWorker.Start()

	// 🔥 worker для чтения ответов
	readMessagesWorker := worker.NewReadMessagesWorker(bhServ, time.Minute)
	readMessagesWorker.Start()

	aHand := handlers.NewAudioHandler(aServ, fServ)

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.Middleware())
	api.HandleFunc("/audio", aHand.Get).Methods("GET")
	api.HandleFunc("/audio", aHand.Add).Methods("POST")
	api.HandleFunc("/audio/{id}", aHand.GetByID).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)

	fmt.Printf("🚀 Сервер запущен на порту %s\n", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, corsHandler))
}
