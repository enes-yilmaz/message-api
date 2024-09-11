package types

type (
	GetSentMessagesResponse struct {
		Items      []Message `json:"items"`
		TotalCount int       `json:"totalCount"`
		Limit      int       `json:"limit"`
		Offset     int       `json:"offset"`
	}
)
