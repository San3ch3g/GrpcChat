package rpc

import (
	pb "ModuleForChat/internal/generated/proto"
	"ModuleForChat/internal/pkg/models"
	"ModuleForChat/internal/pkg/services"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if in.Nickname == "" || in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "nickname or password is empty")
	}

	user := models.User{
		Nickname: in.Nickname,
		Password: in.Password,
	}
	token, err := services.GenerateUserToken(*s.cfg, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	userId, err := s.storage.RegisterUser(&user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.RegisterResponse{UserId: userId, Token: token}, nil
}
