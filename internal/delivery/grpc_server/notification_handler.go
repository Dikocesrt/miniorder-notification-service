package grpc_server

import (
	"context"
	"log/slog"

	"ordermini-notification-service/internal/domain"
	"ordermini-notification-service/internal/interfaces"
	pb "ordermini-notification-service/proto"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	orderUsecase interfaces.IOrderUsecase
	logger       *slog.Logger
}

func NewNotificationHandler(logger *slog.Logger, orderUsecase interfaces.IOrderUsecase) *NotificationHandler {
	return &NotificationHandler{
		orderUsecase: orderUsecase,
		logger:       logger,
	}
}

func (h *NotificationHandler) SendSuccessEmail(ctx context.Context, req *pb.SendSuccessEmailRequest) (*pb.SendSuccessEmailResponse, error) {
	h.logger.InfoContext(ctx, "GRPC INCOMING REQUEST: SendSuccessEmail",
		slog.String("trace_id", req.GetTraceId()),
		slog.String("order_id", req.GetOrderId()),
	)

	orderMsg := &domain.OrderKafkaMessage{
		OrderID:       req.GetOrderId(),
		CustomerEmail: req.GetCustomerEmail(),
		Amount:        req.GetAmount(),
	}

	response := h.orderUsecase.SendSuccessPaymentEmail(ctx, orderMsg)

	if response.Status != 200 {
		h.logger.ErrorContext(ctx, "grpc usecase execution failed",
			slog.String("trace_id", req.GetTraceId()),
			slog.Int("status", response.Status),
		)

		return &pb.SendSuccessEmailResponse{
			Success: false,
			Message: response.Body.(string),
		}, nil
	}

	return &pb.SendSuccessEmailResponse{
		Success: true,
		Message: "Email Notification dispatched successfully",
	}, nil
}
