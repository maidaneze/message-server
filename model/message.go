package model

import "time"

const (
	MAX_TEXT_FIELD_SIZE     = 1024
	MAX_USERNAME_FIELD_SIZE = 32
	MAX_PASSWORD_FIELD_SIZE = 32
)

type MessageResponse struct {
	Id        int64          `json:"id"`
	Timestamp time.Time      `json:"timestamp"`
	Sender    int64          `json:"sender"`
	Recipient int64          `json:"recipient"`
	Content   MessageContent `json:"content"`
}

type MessageContent interface{}

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type ImageContent struct {
	Type   string `json:"type"`
	Url    string `json:"url"`
	Height int64  `json:"height"`
	Width  int64  `json:"width"`
}
type VideoContent struct {
	Type   string `json:"type"`
	Url    string `json:"url"`
	Source string `json:"source"`
}

type MessageDTO struct {
	MessageId   int64
	RecipientId int64
	SenderId    int64
	Timestamp   time.Time
	Type        string
	Text        string
	Url         string
	Height      int64
	Width       int64
	Source      string
}
