package main

import (
	"analize_audio/internal/config"
	"analize_audio/internal/handlers"
	"analize_audio/internal/middleware"
	"analize_audio/internal/repositories"
	"analize_audio/internal/service"
	"analize_audio/pkg/db/gorm"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	cfg := config.Load()

	DB := gorm.InitDB(cfg.DB.Host, cfg.DB.Name, cfg.DB.User, cfg.DB.Password, cfg.DB.Port)
	repositories.Migrate(DB)

	arep := repositories.NewAudioRepository(DB)

	fserv := service.NewFileService(cfg.App.Dir)
	aser := service.NewAudioService(arep)
	//cserv := service.NewConverterService(fserv)
	//cserv.ProcessConvertFiles()
	////_ = cserv
	//
	//bhserv := service.NewBotHubService(cfg.Bot.URL, cfg.Bot.Token, fserv)
	//bhserv.ProcessFileTranscribe()

	ahand := handlers.NewAudioHandler(aser, fserv)

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.Middleware())
	api.HandleFunc("/audio", ahand.Get).Methods("GET")
	api.HandleFunc("/audio", ahand.Add).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(r)

	fmt.Printf("🚀 Сервер запущен на порту %s\n", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, corsHandler))
}
