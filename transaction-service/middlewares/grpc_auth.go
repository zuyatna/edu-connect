package middlewares

import (
	"context"
	"strings"

	"transaction-service/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ContextKey string

const (
	UserIDKey ContextKey = "id"
	EmailKey  ContextKey = "email"
)

func AuthGRPCInterceptor2(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	newCtx := ctx

	if emails, ok := md["email"]; ok && len(emails) > 0 && emails[0] != "" {
		newCtx = context.WithValue(newCtx, EmailKey, emails[0])

		if userIDs, ok := md["id"]; ok && len(userIDs) > 0 && userIDs[0] != "" {
			newCtx = context.WithValue(newCtx, UserIDKey, userIDs[0])
		}
	} else {
		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
		}

		tokenParts := strings.Split(authHeader[0], " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
		}

		tokenString := tokenParts[1]
		userID, email, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		newCtx = context.WithValue(newCtx, EmailKey, email)

		if userID != "" {
			newCtx = context.WithValue(newCtx, UserIDKey, userID)
		}
	}

	return handler(newCtx, req)
}

// func AuthGRPCInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
// 	}

// 	authHeader, ok := md["authorization"]
// 	if !ok || len(authHeader) == 0 {
// 		return nil, status.Errorf(codes.Unauthenticated, "missing authorization header")
// 	}

// 	tokenParts := strings.Split(authHeader[0], " ")
// 	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 		return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
// 	}

// 	tokenString := tokenParts[1]
// 	claims, err := utils.ValidateToken(tokenString)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
// 	}

// 	userID := (*claims)["user_id"].(string)
// 	newCtx := context.WithValue(ctx, UserIDKey, userID)

// 	return handler(newCtx, req)
// }
