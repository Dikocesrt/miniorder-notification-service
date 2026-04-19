package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"ordermini-notification-service/internal/domain"
	"ordermini-notification-service/internal/interfaces"
	"ordermini-notification-service/pkg/common"
	"ordermini-notification-service/pkg/utils"
)

type OrderUsecase struct {
	svcProperties domain.ServiceProperties
	logger        *slog.Logger
}

func NewOrderUsecase(svcProperties domain.ServiceProperties, logger *slog.Logger) interfaces.IOrderUsecase {
	return &OrderUsecase{
		svcProperties: svcProperties,
		logger:        logger,
	}
}

func (u *OrderUsecase) SendSuccessPaymentEmail(ctx context.Context, order *domain.OrderKafkaMessage) common.Response[any] {
	headers := map[string]string{}

	if err := utils.ValidateStruct(order); err != nil {
		u.logger.WarnContext(ctx, "validation failed for order kafka message",
			slog.String("error", err.Error()),
		)
		return common.Response[any]{
			Status: http.StatusBadRequest,
			Header: headers,
			Body:   fmt.Sprintf("validation failed: %s", err.Error()),
		}
	}

	err := utils.SendSuccessPaymentEmail(ctx, u.logger, u.svcProperties, order.OrderID, order.CustomerEmail)
	if err != nil {
		u.logger.ErrorContext(ctx, "failed to send success payment email",
			slog.String("order_id", order.OrderID),
			slog.String("error", err.Error()),
		)
		return common.Response[any]{
			Status: http.StatusInternalServerError,
			Header: headers,
			Body:   "failed to send notification email",
		}
	}

	return common.Response[any]{
		Status: http.StatusOK,
		Header: headers,
		Body:   "notification email sent successfully",
	}
}
