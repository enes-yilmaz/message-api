package message

import (
	"MESSAGEAPI/src/common/postgresql"
	"MESSAGEAPI/src/config"
	"MESSAGEAPI/src/docs"
	messageHandler "MESSAGEAPI/src/internal/messages/handlers"
	"MESSAGEAPI/src/internal/messages/service"
	messageRepository "MESSAGEAPI/src/internal/messages/storages/postgresql"
	fmMiddleware "MESSAGEAPI/src/pkg/middlewares"
	"MESSAGEAPI/src/pkg/utils"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
	"time"
)

func Execute(env string) {
	cfg := config.Config(env)
	config.SetConfig(cfg)

	logrus.Info("MESSAGE API running on \"" + env + "\" environment.")
	logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.ANSIC})

	docs.SwaggerInfo.Host = utils.GetSwagHostEnv()
	e := echo.New()
	e.HideBanner = true
	baseGroup := e.Group("/message-api") // api routing

	baseGroup.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))
	baseGroup.Use(fmMiddleware.PanicExceptionHandling())

	dbPool := postgresql.GetConnectionPool(context.Background(), cfg.PostgresConfig)

	messageRepository := messageRepository.NewMessageRepository(dbPool)

	messageService := service.NewMessageService(messageRepository)

	go messageService.StartSendingMessage()

	messageHandler.NewHandler(baseGroup, messageService)

	baseGroup.GET("/swagger/*", echoSwagger.WrapHandler)

	log.Fatal(e.Start(":4001"))
}
