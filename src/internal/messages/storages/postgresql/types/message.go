package types

import "github.com/google/uuid"

type Message struct {
	ID                   uuid.UUID `json:"id"`
	Content              string    `json:"content"`
	RecipientPhoneNumber string    `json:"recipient_phone_number"`
	Sent                 bool      `json:"sent"`
}
