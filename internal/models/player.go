package models

import "time"

type Player struct {
	ID          int       `json:"id"`
	Plyname     string    `json:"plyname"`
	College     string    `json:"college"`
	Age         string    `json:"age"`
	Height      string    `json:"height"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	//Position       *Position `json:"position,omitempty"`

}
