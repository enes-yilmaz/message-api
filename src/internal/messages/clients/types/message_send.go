package types

type MessageSendRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}
