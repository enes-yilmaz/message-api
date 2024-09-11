package errors

import (
	"MESSAGEAPI/src/pkg/errors"
	"net/http"
)

const repoOp = "messages.repo"
const handlerOp = "messages.handler"
const messageSendClientOp = "message.send.client"

var (
	FailedToBindError   = errors.New(messageSendClientOp, "Failed to bind data from client", 10001, http.StatusInternalServerError)
	SendMessageError    = errors.New(messageSendClientOp, "SendMessageError", 10002, http.StatusInternalServerError)
	JsonUnmarshalFailed = errors.New(handlerOp, "Json unmarshal error", 10003, http.StatusInternalServerError)
)
