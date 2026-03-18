package model

import "time"

type Audio struct {
	ID                      string    `json:"id"`
	SrcFullFilePath         string    `json:"src_file_path"`
	SrcFileName             string    `json:"src_file_name"`
	NeedConvert             bool      `json:"need_convert"`
	AudioFullFilePath       *string   `json:"audio_full_file_path"`
	HasAudio                bool      `json:"has_audio"`
	Transcribed             *string   `json:"transcribed"`
	TranscribedFullFilePath *string   `json:"transcribed_full_file_path"`
	HasTranscribed          bool      `json:"has_transcribed"`
	BotHubMessageID         *string   `json:"bot_hub_message_id"`
	IsMessageSend           bool      `json:"is_message_send"`
	MessageTextFullFilePath *string   `json:"message_text_full_file_path"`
	IsMessageRead           bool      `json:"message_is_read"`
	CreatedAt               time.Time `json:"created_at"`
}
