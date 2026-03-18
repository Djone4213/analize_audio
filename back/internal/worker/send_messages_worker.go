package worker

import (
	"log"
	"time"

	"analize_audio/internal/service"
)

type SendMessagesWorker struct {
	botHubService *service.BotHubService
	interval      time.Duration
	isRunning     bool
}

func NewSendMessagesWorker(botHubService *service.BotHubService, interval time.Duration) *SendMessagesWorker {
	return &SendMessagesWorker{
		botHubService: botHubService,
		interval:      interval,
	}
}

func (w *SendMessagesWorker) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		log.Println("🧠 Send messages worker started")

		for {
			select {
			case <-ticker.C:
				if w.isRunning {
					continue
				}

				w.isRunning = true

				w.botHubService.ProcessSendMessages()

				w.isRunning = false
			}
		}
	}()
}
