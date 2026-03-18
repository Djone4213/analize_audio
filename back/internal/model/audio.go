package model

type Audio struct {
	ID                string  `json:"id"`
	SrcFullFilePath   string  `json:"src_file_path"`
	SrcFileName       string  `json:"src_file_name"`
	NeedConvert       bool    `json:"need_convert"`
	AudioFullFilePath *string `json:"audio_full_file_path"`
	HasAudio          bool    `json:"has_audio"`
	Transcribed       *string `json:"transcribed"`
	HasTranscribed    bool    `json:"has_transcribed"`
}
