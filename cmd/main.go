package main

import (
	"analize_audio/internal/config"
	"analize_audio/internal/service"
)

func main() {
	cfg := config.Load()

	//dir := "D:\\Work\\Projects\\analize_audio\\public\\"

	fserv := service.NewFileService(cfg.App.Dir)
	cserv := service.NewConverterService(fserv)
	cserv.ProcessConvertFiles()
	//_ = cserv

	bhserv := service.NewBotHubService(cfg.Bot.URL, cfg.Bot.Token, fserv)
	bhserv.ProcessFileTranscribe()
}
