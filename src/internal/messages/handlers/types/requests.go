package types

type (
	GetSentMessagesRequest struct {
		Limit          int    `json:"limit"`
		Offset         int    `json:"offset"`
		OrderBy        string `json:"orderBy"`
		OrderDirection string `json:"orderDirection"`
		IsCount        bool   `json:"isCount"`
	}

	MessageAutomationActionRequest struct {
		Action string `json:"action"`
	}
)
