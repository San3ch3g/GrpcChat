package main

import (
	"ModuleForChat/internal/di"
	pb "ModuleForChat/internal/generated/proto"
	"ModuleForChat/internal/pkg/rpc"
	"ModuleForChat/internal/utils/config"
	"fmt"
	"log"
)

func main() {
	cfg := config.NewConfig()
	cfg.InitENV()

	container := di.New(cfg)
	db := container.GetDB()

	postgresDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get database connection: %v", err))
	}
	if err := postgresDB.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}
	grpcServer := container.GetGRPCServer()

	storage := container.GetSQLStorage()
	chatServer := rpc.NewServer(storage, cfg)
	pb.RegisterChatServiceServer(grpcServer, chatServer)
	lis := container.GetNetListener()
	fmt.Println("Server started")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
