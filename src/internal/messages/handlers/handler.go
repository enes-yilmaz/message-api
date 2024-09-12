package handlers

import (
	"MESSAGEAPI/src/internal/messages/errors"
	"MESSAGEAPI/src/internal/messages/handlers/types"
	"MESSAGEAPI/src/internal/messages/helpers"
	"MESSAGEAPI/src/internal/messages/service"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	messageService service.IMessageService
}

func NewHandler(g *echo.Group, messageService service.IMessageService) *Handler {

	h := &Handler{messageService: messageService}

	g = g.Group("/messages")

	g.GET("", h.GetSentMessages)
	g.POST("/automation", h.MessageAutomationAction)

	return h
}

// GetSentMessages
// @Summary  Get sent messages
// @Tags     Messages
// @Accept   json
// @Produce  json
// @Success  200				""
// @Failure  500            	"Internal Error"
// @Router   /messages [get]
func (h Handler) GetSentMessages(c echo.Context) error {

	messages, err := h.messageService.GetSentMessages()

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, types.ToSentMessagesResponse(messages))

}

// MessageAutomationAction
// @Summary  Take action for message automation
// @Tags     Messages
// @Accept   json
// @Produce  json
// @Param    RequestBody  body   types.MessageAutomationActionRequest  true   " "
// @Success  200
// @Failure  400 "Bad Request"
// @Failure  500 "Internal Error"
// @Router   /messages/automation [post]
func (h Handler) MessageAutomationAction(c echo.Context) error {
	var err error

	body, err := io.ReadAll(c.Request().Body)
	req := types.MessageAutomationActionRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(errors.JsonUnmarshalFailed)
	}

	if err = helpers.ValidateMessageAutomationActionRequest(req); err != nil {
		panic(err)
	}

	if strings.TrimSpace(req.Action) == "start" {
		h.messageService.StartSendingMessage()
		return c.NoContent(http.StatusOK)
	} else {
		h.messageService.StopSending()
		return c.NoContent(http.StatusOK)
	}

}
