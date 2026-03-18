package worker

import (
	"log"
	"time"

	"analize_audio/internal/service"
)

type ReadMessagesWorker struct {
	botHubService *service.BotHubService
	interval      time.Duration
	isRunning     bool
}

func NewReadMessagesWorker(botHubService *service.BotHubService, interval time.Duration) *ReadMessagesWorker {
	return &ReadMessagesWorker{
		botHubService: botHubService,
		interval:      interval,
	}
}

func (w *ReadMessagesWorker) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		log.Println("🧠 Read messages worker started")

		for {
			select {
			case <-ticker.C:
				if w.isRunning {
					continue
				}

				w.isRunning = true

				w.botHubService.ProcessReadMessages()

				w.isRunning = false
			}
		}
	}()
}
