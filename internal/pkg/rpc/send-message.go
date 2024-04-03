package rpc

import (
	pb "ModuleForChat/internal/generated/proto"
	"ModuleForChat/internal/pkg/models"
	"ModuleForChat/internal/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func (s *Server) SendMessage(stream pb.ChatService_SendMessageServer) error {
	if err := services.Middleware(stream.Context(), s.cfg); err != nil {
		return err
	}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("Error receiving message from client: %v", err)
			return err
		}

		message := models.Message{
			RoomId:   req.RoomId,
			Content:  req.Content,
			SenderId: req.SenderId,
			Date:     time.Now(),
		}

		messageID, err := s.storage.SaveMessage(&message)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			return status.Error(codes.Internal, "failed to save message")
		}

		if err := stream.SendMsg(&pb.SendMessageResponse{MessageId: messageID}); err != nil {
			log.Printf("Error sending response to client: %v", err)
			return err
		}

		s.broadcastToRoom(req.RoomId, &message)
	}
}

func (s *Server) broadcastToRoom(roomID uint32, message *models.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	streams, ok := s.roomStreams[roomID]
	if !ok {
		return
	}

	for _, stream := range streams {
		if err := stream.Send(&pb.GetMessagesResponse{
			RoomId:    message.RoomId,
			Content:   message.Content,
			Timestamp: message.Date.Format(time.RFC3339),
			SenderId:  message.SenderId,
		}); err != nil {
			log.Printf("Error broadcasting message to client: %v", err)
			streams := s.roomStreams[roomID]
			for i, s := range streams {
				if s == stream {
					streams = append(streams[:i], streams[i+1:]...)
					break
				}
			}
			s.roomStreams[roomID] = streams
		}
	}
}
