package services

import (
	"ModuleForChat/internal/utils/config"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Middleware(ctx context.Context, cfg *config.Config) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	authHeaders := md.Get("Authorization")
	if len(authHeaders) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	authToken := authHeaders[0]
	if authToken == "" {
		return status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	_, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	return nil
}
