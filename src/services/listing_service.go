package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "ebayclone-grpc/proto"
	"ebayclone-grpc/src/storage"
)

type ListingService struct {
	pb.UnimplementedListingServiceServer
	storage storage.Storage
}

func NewListingService(storage storage.Storage) *ListingService {
	return &ListingService{storage: storage}
}

func (s *ListingService) GetListings(ctx context.Context, req *pb.ListingsRequest) (*pb.ListingsResponse, error) {
	listings, err := s.storage.GetListings(req.Search, req.PriceMin, req.PriceMax)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get listings")
	}

	return &pb.ListingsResponse{Listings: listings}, nil
}

func (s *ListingService) CreateListing(ctx context.Context, req *pb.ListingCreate) (*pb.Listing, error) {
	// Validate required fields
	if req.Title == "" || req.Description == "" || req.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Title, description, and price are required")
	}

	// Convert image bytes to base64 strings for storage
	var imageStrings []string
	for i, imageBytes := range req.Images {
		if len(imageBytes) > 5*1024*1024 { // 5MB limit
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Image %d exceeds 5MB limit", i+1))
		}
		imageStrings = append(imageStrings, base64.StdEncoding.EncodeToString(imageBytes))
	}

	if len(imageStrings) > 5 {
		return nil, status.Error(codes.InvalidArgument, "Maximum 5 images allowed")
	}

	listing := &pb.Listing{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		Condition:   req.Condition,
		Location:    req.Location,
		Images:      imageStrings,
		UserId:      getUserIDFromContext(ctx), // Extract from JWT token
	}

	err := s.storage.CreateListing(listing)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create listing")
	}

	return listing, nil
}

func (s *ListingService) GetListing(ctx context.Context, req *pb.GetListingRequest) (*pb.Listing, error) {
	listing, err := s.storage.GetListing(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "Listing not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get listing")
	}
	return listing, nil
}

func (s *ListingService) UpdateListing(ctx context.Context, req *pb.UpdateListingRequest) (*pb.Listing, error) {
	// Get existing listing
	existing, err := s.storage.GetListing(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "Listing not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get listing")
	}

	// Update fields if provided
	updated := &pb.Listing{
		Id:          existing.Id,
		Title:       existing.Title,
		Description: existing.Description,
		Price:       existing.Price,
		Category:    existing.Category,
		Condition:   existing.Condition,
		Location:    existing.Location,
		Images:      existing.Images,
		UserId:      existing.UserId,
		CreatedAt:   existing.CreatedAt,
	}

	if req.Listing.Title != "" {
		updated.Title = req.Listing.Title
	}
	if req.Listing.Description != "" {
		updated.Description = req.Listing.Description
	}
	if req.Listing.Price > 0 {
		updated.Price = req.Listing.Price
	}
	if req.Listing.Category != "" {
		updated.Category = req.Listing.Category
	}
	if req.Listing.Condition != "" {
		updated.Condition = req.Listing.Condition
	}
	if req.Listing.Location != "" {
		updated.Location = req.Listing.Location
	}

	updated.UpdatedAt = timestamppb.New(time.Now())

	err = s.storage.UpdateListing(req.Id, updated)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to update listing")
	}

	return updated, nil
}

func (s *ListingService) DeleteListing(ctx context.Context, req *pb.DeleteListingRequest) (*pb.Success, error) {
	err := s.storage.DeleteListing(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "Listing not found")
		}
		return nil, status.Error(codes.Internal, "Failed to delete listing")
	}

	return &pb.Success{Message: "Listing deleted successfully"}, nil
}

// Helper function to extract user ID from JWT token in context
func getUserIDFromContext(ctx context.Context) int32 {
	// In a real implementation, extract from JWT token in metadata
	// For demo purposes, return a default user ID
	return 1
}
