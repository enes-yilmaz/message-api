package types

type (
	Message struct {
		ID                   string `json:"id"`
		Content              string `json:"content"`
		RecipientPhoneNumber string `json:"recipient_phone_number"`
		Sent                 bool   `json:"sent"`
	}
)
