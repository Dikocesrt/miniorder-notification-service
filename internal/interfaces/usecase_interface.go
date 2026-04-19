package interfaces

import (
	"context"
	"ordermini-notification-service/internal/domain"
	"ordermini-notification-service/pkg/common"
)

type IOrderUsecase interface {
	SendSuccessPaymentEmail(ctx context.Context, order *domain.OrderKafkaMessage) common.Response[any]
}
