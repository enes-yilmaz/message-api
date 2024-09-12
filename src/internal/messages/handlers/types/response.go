package types

import "MESSAGEAPI/src/internal/messages/models"

type (
	GetSentMessagesResponse struct {
		Messages   []Message `json:"messages"`
		TotalCount int       `json:"totalCount"`
	}
)

func ToSentMessagesResponse(messages []models.Message) GetSentMessagesResponse {
	var messagesResponseList = make([]Message, 0)
	for _, message := range messages {
		messagesResponseList = append(messagesResponseList, Message{
			Content:              message.Content,
			RecipientPhoneNumber: message.RecipientPhoneNumber,
			Sent:                 message.Sent,
		})
	}
	return GetSentMessagesResponse{
		Messages:   messagesResponseList,
		TotalCount: len(messages),
	}
}
