package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"

	"ordermini-notification-service/internal/config"
	"ordermini-notification-service/internal/delivery/grpc_server"
	"ordermini-notification-service/internal/usecase"
	pb "ordermini-notification-service/proto"

	"google.golang.org/grpc"
)

func RunApplication() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	svcProperties := config.GetServiceProperties()
	fmt.Println(strings.ToUpper(svcProperties.ServiceName) + " SERVICE")
	asBytesJSON, _ := json.Marshal(svcProperties)
	fmt.Printf("SERVICE CONFIG : '%v'\n", string(asBytesJSON))

	grpcServer := grpc.NewServer()

	iOrderUsecase := usecase.NewOrderUsecase(svcProperties, logger)

	notificationHandler := grpc_server.NewNotificationHandler(logger, iOrderUsecase)

	pb.RegisterNotificationServiceServer(grpcServer, notificationHandler)

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(svcProperties.ServiceTCPPort))
	if err != nil {
		logger.Error("failed to open tcp listen for port "+strconv.Itoa(svcProperties.ServiceTCPPort), "error", err)
		os.Exit(1)
	}

	logger.Info("GRPC Server Notification-Service nyala dan bersiaga di port :" + strconv.Itoa(svcProperties.ServiceTCPPort))

	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("failed to serve grpc", "error", err)
	}
}
