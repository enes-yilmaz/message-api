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
	messageRepository messageRepository.IMessageRepository
	messageSendClient *clients.MessageSend
	stopChan          chan struct{}
	wg                sync.WaitGroup
	isSendingRunning  bool
	redisClient       *redis.Client
	sentMessageMap    sync.Map
}

func NewMessageService(messageRepository messageRepository.IMessageRepository) IMessageService {

	messageClient := clients.NewMessageSendClient()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.GetConfigs().RedisConfig.Addr,
		Password: cfg.GetConfigs().RedisConfig.Password,
		DB:       cfg.GetConfigs().RedisConfig.DB,
	})

	return &MessageService{
		messageRepository: messageRepository,
		messageSendClient: messageClient,
		stopChan:          make(chan struct{}),
		isSendingRunning:  false,
		redisClient:       redisClient,
	}
}

func (ms *MessageService) GetSentMessages() ([]models.Message, error) {

	messages, err := ms.messageRepository.GetSentMessages()
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (ms *MessageService) StartSendingMessage() {
	if ms.isSendingRunning {
		log.Info("Message sending is already running")
		return
	}

	ms.isSendingRunning = true
	duration := cfg.GetConfigs().MessageSendDuration
	ticker := time.NewTicker(duration * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ms.wg.Add(1)
			go func() {
				defer ms.wg.Done()
				log.Info("Sending pending messages")
				if err := ms.sendPendingMessages(); err != nil {
					log.Error("Error while sending messages: %v", err)
				}
			}()
		case <-ms.stopChan:
			log.Info("Message sending stopped")
			ms.wg.Wait()
			ms.isSendingRunning = false
			return
		}
	}
}

func (ms *MessageService) StopSending() {
	if !ms.isSendingRunning {
		return
	}

	close(ms.stopChan)
	ms.wg.Wait()
	ms.sentMessageMap.Range(func(key, _ any) bool {
		ms.sentMessageMap.Delete(key)
		return true
	})
	ms.isSendingRunning = false
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
			// Check if message has already been sent (concurrent-safe)
			if _, exists := ms.sentMessageMap.Load(msg.ID); exists {
				return nil
			}

			ms.sentMessageMap.Store(msg.ID, struct{}{}) // Mark as sent

			tx, err := ms.messageRepository.BeginTransaction()
			if err != nil {
				return err
			}
			defer tx.Rollback(ctx)

			if err := ms.sendMessageToPhone(msg); err != nil {
				return err
			}

			if err := ms.messageRepository.MarkMessageAsSent(msg.ID, tx); err != nil {
				return err
			}

			if err := tx.Commit(ctx); err != nil {
				return err
			}
			// Remove from cache after successful sending
			ms.sentMessageMap.Delete(msg.ID)
			return nil
		})
	}

	return g.Wait()
}

func (ms *MessageService) sendMessageToPhone(message models.Message) error {
	if err := ms.messageSendClient.SendMessage(types.MessageSendRequest{
		Content: message.Content,
		To:      message.RecipientPhoneNumber,
	}); err != nil {
		return err
	}

	// Cache message ID and sending time in Redis
	sendingTime := time.Now().Unix()
	if err := ms.redisClient.Set(context.Background(), message.ID, sendingTime, 5*time.Hour).Err(); err != nil {
		log.Errorf("Failed to cache message ID to Redis: %v", err)
	}

	return nil
}
