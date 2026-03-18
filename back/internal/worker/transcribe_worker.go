package worker

import (
	"log"
	"time"

	"analize_audio/internal/service"
)

type TranscribeWorker struct {
	botHubService *service.BotHubService
	interval      time.Duration
	isRunning     bool
}

func NewTranscribeWorker(botHubService *service.BotHubService, interval time.Duration) *TranscribeWorker {
	return &TranscribeWorker{
		botHubService: botHubService,
		interval:      interval,
	}
}

func (w *TranscribeWorker) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		log.Println("🧠 Transcribe worker started")

		for {
			select {
			case <-ticker.C:
				if w.isRunning {
					continue
				}

				w.isRunning = true

				w.botHubService.ProcessFileTranscribe()

				w.isRunning = false
			}
		}
	}()
}
