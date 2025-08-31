package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "ebayclone-grpc/proto"
	"ebayclone-grpc/src/services"
	"ebayclone-grpc/src/storage"
)

func main() {
	// Initialize storage
	store := storage.NewInMemoryStorage()

	// Create gRPC server
	s := grpc.NewServer()

	// Register services
	pb.RegisterUserServiceServer(s, services.NewUserService(store))
	pb.RegisterSessionServiceServer(s, services.NewSessionService(store))
	pb.RegisterListingServiceServer(s, services.NewListingService(store))
	pb.RegisterOrderServiceServer(s, services.NewOrderService(store))

	// Enable reflection for testing
	reflection.Register(s)

	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("gRPC server starting on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
