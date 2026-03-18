package model

type ChatModel struct {
	ModelId   string `json:"modelId"`
	Name      string `json:"name"`
	Highlight string `json:"highlight"`
	Platform  string `json:"platform"`
}
