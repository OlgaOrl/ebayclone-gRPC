package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "ebayclone-grpc/proto"
	"ebayclone-grpc/src/storage"
)

type SessionService struct {
	pb.UnimplementedSessionServiceServer
	storage   storage.Storage
	jwtSecret []byte
}

func NewSessionService(storage storage.Storage) *SessionService {
	return &SessionService{
		storage:   storage,
		jwtSecret: []byte("your-secret-key"), // In production, use environment variable
	}
}

func (s *SessionService) Login(ctx context.Context, req *pb.UserLogin) (*pb.LoginResponse, error) {
	// Validate required fields
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Email and password are required")
	}

	// Get user by email
	user, err := s.storage.GetUserByEmail(req.Email)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	// Verify password
	hashedPassword := hashPassword(req.Password)
	if memStorage, ok := s.storage.(*storage.InMemoryStorage); ok {
		storedPassword, exists := memStorage.GetUserPassword(user.Id)
		if !exists || storedPassword != hashedPassword {
			return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
		}
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}

func (s *SessionService) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	// In a real implementation, you might invalidate the token in a blacklist
	// For this demo, we just return success
	return &emptypb.Empty{}, nil
}
