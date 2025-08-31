package services

import (
	"context"
	"math"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "ebayclone-grpc/proto"
	"ebayclone-grpc/src/storage"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	storage storage.Storage
}

func NewOrderService(storage storage.Storage) *OrderService {
	return &OrderService{storage: storage}
}

func (s *OrderService) GetOrders(ctx context.Context, req *pb.OrdersRequest) (*pb.OrdersResponse, error) {
	orders, total, err := s.storage.GetOrders(req.UserId, req.Status, req.Page, req.Limit)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get orders")
	}

	// Calculate pagination
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	pages := int32(math.Ceil(float64(total) / float64(limit)))

	pagination := &pb.Pagination{
		Total: total,
		Pages: pages,
	}

	return &pb.OrdersResponse{
		Orders:     orders,
		Pagination: pagination,
	}, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.OrderCreate) (*pb.Order, error) {
	// Validate required fields
	if req.ListingId <= 0 || req.Quantity <= 0 || req.ShippingAddress == nil {
		return nil, status.Error(codes.InvalidArgument, "ListingId, quantity, and shipping address are required")
	}

	// Validate shipping address
	addr := req.ShippingAddress
	if addr.Street == "" || addr.City == "" || addr.Country == "" {
		return nil, status.Error(codes.InvalidArgument, "Street, city, and country are required in shipping address")
	}

	// Get listing to calculate total price
	listing, err := s.storage.GetListing(req.ListingId)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "Listing not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get listing")
	}

	totalPrice := listing.Price * float64(req.Quantity)

	order := &pb.Order{
		UserId:          getUserIDFromContext(ctx),
		ListingId:       req.ListingId,
		Quantity:        req.Quantity,
		TotalPrice:      totalPrice,
		Status:          "pending",
		ShippingAddress: req.ShippingAddress,
		BuyerNotes:      req.BuyerNotes,
	}

	err = s.storage.CreateOrder(order)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create order")
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := s.storage.GetOrder(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "Order not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get order")
	}
	return order, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.Order, error) {
	// Get existing order
	existing, err := s.storage.GetOrder(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "Order not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get order")
	}

	// Update fields if provided
	updated := &pb.Order{
		Id:              existing.Id,
		UserId:          existing.UserId,
		ListingId:       existing.ListingId,
		Quantity:        existing.Quantity,
		TotalPrice:      existing.TotalPrice,
		Status:          existing.Status,
		ShippingAddress: existing.ShippingAddress,
		BuyerNotes:      existing.BuyerNotes,
		CreatedAt:       existing.CreatedAt,
		CancelledAt:     existing.CancelledAt,
		CancelReason:    existing.CancelReason,
	}

	if req.Order.UserId > 0 {
		updated.UserId = req.Order.UserId
	}
	if req.Order.ListingId > 0 {
		updated.ListingId = req.Order.ListingId
	}
	if req.Order.Quantity > 0 {
		updated.Quantity = req.Order.Quantity
	}
	if req.Order.TotalPrice > 0 {
		updated.TotalPrice = req.Order.TotalPrice
	}

	updated.UpdatedAt = timestamppb.New(time.Now())

	err = s.storage.UpdateOrder(req.Id, updated)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to update order")
	}

	return updated, nil
}

func (s *OrderService) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.Success, error) {
	err := s.storage.DeleteOrder(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "Order not found")
		}
		return nil, status.Error(codes.Internal, "Failed to delete order")
	}

	return &pb.Success{Message: "Order deleted successfully"}, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	// Get existing order
	existing, err := s.storage.GetOrder(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "Order not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get order")
	}

	// Update order status to cancelled
	updated := &pb.Order{
		Id:              existing.Id,
		UserId:          existing.UserId,
		ListingId:       existing.ListingId,
		Quantity:        existing.Quantity,
		TotalPrice:      existing.TotalPrice,
		Status:          "cancelled",
		ShippingAddress: existing.ShippingAddress,
		BuyerNotes:      existing.BuyerNotes,
		CreatedAt:       existing.CreatedAt,
		UpdatedAt:       timestamppb.New(time.Now()),
		CancelledAt:     timestamppb.New(time.Now()),
		CancelReason:    req.CancelReason,
	}

	err = s.storage.UpdateOrder(req.Id, updated)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to cancel order")
	}

	return &pb.CancelOrderResponse{
		Message: "Order cancelled successfully",
		Order:   updated,
	}, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {
	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"shipped":   true,
		"delivered": true,
	}

	if !validStatuses[req.Status] {
		return nil, status.Error(codes.InvalidArgument, "Invalid status. Must be: pending, confirmed, shipped, or delivered")
	}

	// Get existing order
	existing, err := s.storage.GetOrder(req.Id)
	if err != nil {
		if _, ok := err.(*storage.NotFoundError); ok {
			return nil, status.Error(codes.NotFound, "Order not found")
		}
		return nil, status.Error(codes.Internal, "Failed to get order")
	}

	// Update order status
	updated := &pb.Order{
		Id:              existing.Id,
		UserId:          existing.UserId,
		ListingId:       existing.ListingId,
		Quantity:        existing.Quantity,
		TotalPrice:      existing.TotalPrice,
		Status:          req.Status,
		ShippingAddress: existing.ShippingAddress,
		BuyerNotes:      existing.BuyerNotes,
		CreatedAt:       existing.CreatedAt,
		UpdatedAt:       timestamppb.New(time.Now()),
		CancelledAt:     existing.CancelledAt,
		CancelReason:    existing.CancelReason,
	}

	err = s.storage.UpdateOrder(req.Id, updated)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to update order status")
	}

	return updated, nil
}
