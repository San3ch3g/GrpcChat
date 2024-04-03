package rpc

import (
	pb "ModuleForChat/internal/generated/proto"
	"ModuleForChat/internal/pkg/models"
	"ModuleForChat/internal/pkg/storage/pg"
	"ModuleForChat/internal/utils/config"
	"sync"
)

type Server struct {
	*pb.UnimplementedChatServiceServer
	storage           *pg.Storage
	cfg               *config.Config
	roomStreams       map[uint32][]pb.ChatService_GetMessagesServer
	mu                sync.Mutex
	newMessageChannel chan *models.Message
}

func NewServer(storage *pg.Storage, cfg *config.Config) *Server {
	return &Server{
		storage:           storage,
		cfg:               cfg,
		roomStreams:       make(map[uint32][]pb.ChatService_GetMessagesServer),
		newMessageChannel: make(chan *models.Message),
	}
}
