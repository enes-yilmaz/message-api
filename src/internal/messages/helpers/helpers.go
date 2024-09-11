package helpers

import (
	"MESSAGEAPI/src/internal/messages/handlers/types"
	"MESSAGEAPI/src/pkg/errors"
	"fmt"
	"strings"
)

func ValidateMessageAutomationActionRequest(req types.MessageAutomationActionRequest) error {

	// Action Validation
	validMessageActions := map[string]bool{
		"start": true,
		"stop":  true,
	}
	if _, ok := validMessageActions[strings.TrimSpace(req.Action)]; !ok {
		return errors.ValidatorError.WrapDesc(fmt.Sprintf("Invalid Action value:%s. Valid Action values are start,stop.", req.Action))
	}
	return nil
}
