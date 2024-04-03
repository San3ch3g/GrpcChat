package rpc

import (
	pb "ModuleForChat/internal/generated/proto"
	"ModuleForChat/internal/pkg/models"
	"ModuleForChat/internal/pkg/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

func (s *Server) GetMessages(req *pb.GetMessagesRequest, stream pb.ChatService_GetMessagesServer) error {
	if err := services.Middleware(stream.Context(), s.cfg); err != nil {
		return err
	}

	roomID := req.RoomId
	s.mu.Lock()
	s.roomStreams[roomID] = append(s.roomStreams[roomID], stream)
	s.mu.Unlock()

	messages, err := s.storage.GetMessagesByRoomID(roomID)
	if err != nil {
		return status.Error(codes.Internal, "failed to get messages")
	}

	for _, message := range messages {
		response := &pb.GetMessagesResponse{
			RoomId:    message.RoomId,
			Content:   message.Content,
			Timestamp: message.Date.Format(time.RFC3339),
			SenderId:  message.SenderId,
		}

		if err := stream.Send(response); err != nil {
			return err
		}
	}

	for {
		select {
		case <-stream.Context().Done():
			s.mu.Lock()
			streams := s.roomStreams[roomID]
			for i, s := range streams {
				if s == stream {
					streams = append(streams[:i], streams[i+1:]...)
					break
				}
			}
			s.roomStreams[roomID] = streams
			s.mu.Unlock()
			return nil
		case newMessage := <-s.newMessageChannel:
			if newMessage.RoomId != roomID {
				continue
			}

			response := &pb.GetMessagesResponse{
				RoomId:    newMessage.RoomId,
				Content:   newMessage.Content,
				Timestamp: newMessage.Date.Format(time.RFC3339),
				SenderId:  newMessage.SenderId,
			}

			if err := stream.Send(response); err != nil {
				s.mu.Lock()
				streams := s.roomStreams[roomID]
				for i, s := range streams {
					if s == stream {
						streams = append(streams[:i], streams[i+1:]...)
						break
					}
				}
				s.roomStreams[roomID] = streams
				s.mu.Unlock()
				return err
			}
		}
	}
}

func (s *Server) broadcastNewMessage(roomID uint32, message *models.Message) {
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
		}
	}
}
