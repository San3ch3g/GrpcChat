package rpc

import (
	pb "ModuleForChat/internal/generated/proto"
	"ModuleForChat/internal/pkg/models"
	"ModuleForChat/internal/pkg/services"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	if in.Nickname == "" || in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "nickname or password is empty")
	}

	user := models.User{
		Nickname: in.Nickname,
		Password: in.Password,
	}
	userId, err := s.storage.LoginUser(&user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	user.Id = userId

	token, err := services.GenerateUserToken(*s.cfg, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &pb.LoginResponse{UserId: userId, Token: token}, nil
}
