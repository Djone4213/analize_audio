package worker

import (
	"log"
	"time"

	"analize_audio/internal/service"
)

type ConverterWorker struct {
	converterService *service.ConverterService
	interval         time.Duration
	isRunning        bool
}

func NewConverterWorker(converterService *service.ConverterService, interval time.Duration) *ConverterWorker {
	return &ConverterWorker{
		converterService: converterService,
		interval:         interval,
	}
}

func (w *ConverterWorker) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		log.Println("🧠 Converts worker started")

		for {
			select {
			case <-ticker.C:
				if w.isRunning {
					continue
				}

				w.isRunning = true
				w.converterService.ProcessConvertFiles()
				w.isRunning = false
			}
		}
	}()
}
