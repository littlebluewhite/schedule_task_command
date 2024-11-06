package grpc_task_template

import (
	"context"
	"fmt"
	"github.com/littlebluewhite/schedule_task_command/api"
	"github.com/littlebluewhite/schedule_task_command/util/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func StartGRPCServer(ctx context.Context, mainLog *logrus.Logger,
	Config config.Config, ts TaskServer, d api.Dbs) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", Config.Server.GPort))
	if err != nil {
		mainLog.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register gRPC services
	RegisterTaskTemplateServiceServer(grpcServer, NewTaskTemplateService(ts, d))

	mainLog.Infof("gRPC server started on port %s", Config.Server.GPort)

	go func() {
		<-ctx.Done()
		mainLog.Infoln("Gracefully shutting down gRPC server")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		mainLog.Fatalf("failed to serve: %v", err)
	}
}
