package domain

import (
	"efficient-api/utils/error_utils"
	"strings"
	"time"
)

type Message struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func (m *Message) Validate() error_utils.MessageErr {
	m.Title = strings.TrimSpace(m.Title)
	m.Body = strings.TrimSpace(m.Body)
	if m.Title == "" {
		return error_utils.NewBadRequestError("Please enter a valid title")
	}
	if m.Body == "" {
		return error_utils.NewBadRequestError("Please enter a valid body")
	}
	return nil
}
