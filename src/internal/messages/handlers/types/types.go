package types

type (
	Message struct {
		Content              string `json:"content"`
		RecipientPhoneNumber string `json:"recipient_phone_number"`
		Sent                 bool   `json:"sent"`
	}
)
