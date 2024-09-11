package clients

import (
	cfg "MESSAGEAPI/src/config"
	"MESSAGEAPI/src/internal/messages/clients/types"
	"MESSAGEAPI/src/internal/messages/errors"
	rest "MESSAGEAPI/src/pkg/clients"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type MessageSend struct {
	client *rest.BaseClient
}

func NewMessageSendClient() *MessageSend {
	return &MessageSend{
		client: rest.NewBaseClient(cfg.GetConfigs().MessageSendClient),
	}
}

func (ms MessageSend) SendMessage(request types.MessageSendRequest) error {

	reqBody, err := json.Marshal(request)
	if err != nil {
		return errors.FailedToBindError
	}

	res, err := ms.client.POST("", reqBody)
	if err != nil {
		return errors.SendMessageError
	}

	defer fasthttp.ReleaseResponse(res)

	//if res.StatusCode() != 202 {
	if res.StatusCode() > 200 && res.StatusCode() < 300 {
		return errors.SendMessageError
	}

	return nil

}
