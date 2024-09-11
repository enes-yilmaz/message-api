package service

import (
	cfg "MESSAGEAPI/src/config"
	"MESSAGEAPI/src/internal/messages/clients"
	"MESSAGEAPI/src/internal/messages/clients/types"
	"MESSAGEAPI/src/internal/messages/models"
	messageRepository "MESSAGEAPI/src/internal/messages/storages/postgresql"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

type IMessageService interface {
	GetSentMessages() ([]models.Message, error)
	StartSendingMessage()
	StopSending()
}

type MessageService struct {
	messageRepository    messageRepository.IMessageRepository
	messageSendClient    *clients.MessageSend
	messageMap           map[string]struct{}
	stopChan             chan struct{}
	wg                   sync.WaitGroup
	isSentMessageRunning bool
	redisClient          *redis.Client
}

func NewMessageService(messageRepository messageRepository.IMessageRepository) IMessageService {

	messageClient := clients.NewMessageSendClient()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.GetConfigs().RedisConfig.Addr,
		Password: cfg.GetConfigs().RedisConfig.Password,
		DB:       cfg.GetConfigs().RedisConfig.DB,
	})

	return &MessageService{
		messageRepository:    messageRepository,
		messageSendClient:    messageClient,
		messageMap:           make(map[string]struct{}),
		stopChan:             make(chan struct{}),
		isSentMessageRunning: false,
		redisClient:          redisClient,
	}
}

func (ms *MessageService) GetSentMessages() ([]models.Message, error) {

	messages, err := ms.messageRepository.GetAllSentMessages()
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (ms *MessageService) StartSendingMessage() {
	if ms.isSentMessageRunning {
		log.Info("Message sending is already running")
		return
	}
	ms.isSentMessageRunning = true
	duration := cfg.GetConfigs().MessageSendDuration
	ticker := time.NewTicker(duration * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ms.wg.Add(1)
			func() {
				defer ms.wg.Done()
				log.Info("Sending pending messages")
				err := ms.sendPendingMessages()
				if err != nil {
					log.Error("Error while sending messages: %v", err)
				}
			}()
		case <-ms.stopChan:
			log.Info("Message sending stopped")
			ms.wg.Wait()
			return
		}
	}
}

func (ms *MessageService) StopSending() {
	close(ms.stopChan)
	ms.wg.Wait()
	ms.messageMap = make(map[string]struct{})
	ms.isSentMessageRunning = false
}

func (ms *MessageService) sendPendingMessages() error {
	unsentMessages, err := ms.messageRepository.GetUnSentMessages()
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(context.Background())

	for _, message := range unsentMessages {
		msg := message
		g.Go(func() error {
			if _, exists := ms.messageMap[msg.ID]; exists {
				return nil
			}
			ms.messageMap[msg.ID] = struct{}{}

			tx, err := ms.messageRepository.BeginTransaction()
			if err != nil {
				return err
			}
			defer tx.Rollback(ctx)

			err = ms.sendMessageToPhone(msg)
			if err != nil {
				return err
			}

			err = ms.messageRepository.MarkMessageAsSent(msg.ID, tx)
			if err != nil {
				return err
			}

			if err := tx.Commit(ctx); err != nil {
				return err
			}

			delete(ms.messageMap, msg.ID)
			return nil
		})
	}

	return g.Wait()
}

func (ms *MessageService) sendMessageToPhone(message models.Message) error {
	err := ms.messageSendClient.SendMessage(types.MessageSendRequest{
		Content: message.Content,
		To:      message.RecipientPhoneNumber,
	})
	if err != nil {
		return err
	}

	messageId := message.ID
	sendingTime := time.Now().Unix()
	err = ms.redisClient.Set(context.Background(), messageId, sendingTime, 5*time.Hour).Err()
	if err != nil {
		log.Errorf("Failed to cache message ID to Redis: %v", err)
	}
	return nil
}
