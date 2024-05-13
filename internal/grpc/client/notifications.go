package client

import (
	"context"
	"github.com/SiriusServiceDesk/auth-service/internal/config"
	"github.com/SiriusServiceDesk/auth-service/pkg/logger"
	"github.com/SiriusServiceDesk/gateway-service/pkg/notification_v1"

	"google.golang.org/grpc"
)

type Message struct {
	To           []string
	Data         string
	Type         string
	Subject      string
	TemplateName string
}

func createConnectionToNotification() (notification_v1.NotificationV1Client, error) {
	cfg := config.GetConfig()
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, cfg.NotificationService.Address, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("Failed to connect gRPC notification:%s\n", err.Error())
		return nil, err
	}
	return notification_v1.NewNotificationV1Client(conn), nil
}

func SendMessage(m *Message) (*notification_v1.CreateMessageResponse, error) {
	ctx := context.Background()
	conn, err := createConnectionToNotification()
	if err != nil {
		return nil, err
	}

	message := &notification_v1.CreateMessageRequest{
		Subject:      m.Subject,
		To:           m.To,
		Data:         m.Data,
		TemplateName: m.TemplateName,
		Type:         m.Type,
	}

	response, err := conn.CreateMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	return response, nil
}
