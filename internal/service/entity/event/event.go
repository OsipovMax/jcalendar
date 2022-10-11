package event

import "time"

type Event struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	From      string    `json:"from"`
	Till      string    `json:"till"`
	CreatorID string    `json:"creator_id"`
	Details   string    `json:"details"`
	IsPrivate bool      `json:"is_private"`
	IsRepeat  bool      `json:"is_repeat"`
}
