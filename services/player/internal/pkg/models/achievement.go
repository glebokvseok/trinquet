package models

import "time"

type Achievement struct {
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Description string     `json:"description"`
	IsCompleted bool       `json:"is_completed"`
	CompletedOn *time.Time `json:"completed_on"`
}
