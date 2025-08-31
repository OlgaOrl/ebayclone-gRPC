package services

import (
	"context"
	"crypto/sha256"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "ebayclone-grpc/proto"
	"ebayclone-grpc/src/storage"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	storage storage.Storage
}

func NewUserService(storage storage.Storage) *UserService {
	return &UserService{storage: storage}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.UserCreate) (*pb.User, error) {
	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "All fields are required")
	}

	// Hash password
	hashedPassword := hashPassword(req.Password)

	user := &pb.User{
		Username: req.Username,
		Email:    req.Email,
	}

	err := s.storage.CreateUser(user)
	if err != nil {
		if _, ok := err.(*storage.UserExistsError); ok {
			return nil, status.Error(codes.AlreadyExists, "Email already exists")
		}
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	// Store password separately
	if memStorage, ok := s.storage.(*storage.InMemoryStorage); ok {
		memStorage.SetUserPassword(user.Id, hashedPassword)
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	user, err := s.storage.GetUser(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get user")
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	// Get existing user
	existing, err := s.storage.GetUser(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get user")
	}

	// Update fields if provided
	updated := &pb.User{
		Id:       existing.Id,
		Username: existing.Username,
		Email:    existing.Email,
	}

	if req.User.Username != "" {
		updated.Username = req.User.Username
	}
	if req.User.Email != "" {
		updated.Email = req.User.Email
	}

	err = s.storage.UpdateUser(req.Id, updated)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to update user")
	}

	// Update password if provided
	if req.User.Password != "" {
		hashedPassword := hashPassword(req.User.Password)
		if memStorage, ok := s.storage.(*storage.InMemoryStorage); ok {
			memStorage.SetUserPassword(req.Id, hashedPassword)
		}
	}

	return updated, nil
}

func (s *UserService) ReplaceUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	// Check if user exists
	_, err := s.storage.GetUser(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get user")
	}

	// Validate required fields for replacement
	if req.User.Username == "" || req.User.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "Username and email are required")
	}

	user := &pb.User{
		Id:       req.Id,
		Username: req.User.Username,
		Email:    req.User.Email,
	}

	err = s.storage.UpdateUser(req.Id, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to replace user")
	}

	// Update password if provided
	if req.User.Password != "" {
		hashedPassword := hashPassword(req.User.Password)
		if memStorage, ok := s.storage.(*storage.InMemoryStorage); ok {
			memStorage.SetUserPassword(req.Id, hashedPassword)
		}
	}

	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.storage.DeleteUser(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Error(codes.Internal, "Failed to delete user")
	}
	return &emptypb.Empty{}, nil
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}
