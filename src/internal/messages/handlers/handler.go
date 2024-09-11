package handlers

import (
	"MESSAGEAPI/src/internal/messages/errors"
	"MESSAGEAPI/src/internal/messages/handlers/types"
	"MESSAGEAPI/src/internal/messages/helpers"
	"MESSAGEAPI/src/internal/messages/service"
	"MESSAGEAPI/src/pkg/utils"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
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
// @Param    limit     			query    number    	false " " "10"
// @Param    offset      		query    number		false " " "0"
// @Param    orderDirection  	query    string     false "orderDirection"
// @Param    orderBy         	query    string     false "orderBy"
// @Param    isCount         	query    boolean    false "false"
// @Success  200				""
// @Failure  404 				"Not Found"
// @Failure  500            	"Internal Error"
// @Router   /messages [get]
func (h Handler) GetSentMessages(c echo.Context) error {

	req := types.GetSentMessagesRequest{}

	if limit, offset, err := utils.ValidateLimitOffset(strings.TrimSpace(c.QueryParam("limit")), strings.TrimSpace(c.QueryParam("offset")), 10); err != nil {
		panic(err)
	} else {
		req.Limit, req.Offset = limit, offset
	}
	req.OrderDirection = strings.TrimSpace(c.QueryParam("orderDirection"))
	req.OrderBy = strings.TrimSpace(c.QueryParam("orderBy"))
	req.IsCount, _ = strconv.ParseBool(strings.TrimSpace(c.QueryParam("isCount")))

	messages, err := h.messageService.GetSentMessages()
	if err != nil {
		return err
	}
	return c.JSON(200, messages)

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
