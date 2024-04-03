package rpc

import (
	pb "ModuleForChat/internal/generated/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateRoom(ctx context.Context, in *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {

	if in.RoomName == "" {
		return nil, status.Error(codes.InvalidArgument, "room name is required")
	}

	if s.storage == nil {
		return nil, status.Error(codes.Internal, "storage is not initialized")
	}

	room, err := s.storage.CreateRoom(in.RoomName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &pb.CreateRoomResponse{Room: &pb.Room{Id: room.Id, Name: room.Name}}
	return response, nil
}
