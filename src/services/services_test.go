package services

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "ebayclone-grpc/proto"
	"ebayclone-grpc/src/storage"
)

func TestUserService(t *testing.T) {
	store := storage.NewInMemoryStorage()
	service := NewUserService(store)
	ctx := context.Background()

	// Test CreateUser
	user, err := service.CreateUser(ctx, &pb.UserCreate{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if user.Username != "testuser" || user.Email != "test@example.com" {
		t.Errorf("User data mismatch: got %+v", user)
	}

	// Test GetUser
	retrievedUser, err := service.GetUser(ctx, &pb.GetUserRequest{Id: user.Id})
	if err != nil {
		t.Fatalf("GetUser failed: %v", err)
	}
	if retrievedUser.Id != user.Id {
		t.Errorf("Retrieved user ID mismatch: expected %d, got %d", user.Id, retrievedUser.Id)
	}

	// Test UpdateUser
	updatedUser, err := service.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:   user.Id,
		User: &pb.UserUpdate{Username: "updateduser"},
	})
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}
	if updatedUser.Username != "updateduser" {
		t.Errorf("Username not updated: expected 'updateduser', got '%s'", updatedUser.Username)
	}

	// Test error cases
	_, err = service.GetUser(ctx, &pb.GetUserRequest{Id: 999})
	if status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound error for non-existent user, got: %v", err)
	}

	// Test CreateUser with duplicate email
	_, err = service.CreateUser(ctx, &pb.UserCreate{
		Username: "anotheruser",
		Email:    "test@example.com", // Same email
		Password: "password123",
	})
	if status.Code(err) != codes.AlreadyExists {
		t.Errorf("Expected AlreadyExists error for duplicate email, got: %v", err)
	}
}

func TestSessionService(t *testing.T) {
	store := storage.NewInMemoryStorage()
	userService := NewUserService(store)
	sessionService := NewSessionService(store)
	ctx := context.Background()

	// Create a user first
	_, err := userService.CreateUser(ctx, &pb.UserCreate{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test Login
	loginResp, err := sessionService.Login(ctx, &pb.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if loginResp.Token == "" {
		t.Error("Login should return a token")
	}

	// Test Login with wrong password
	_, err = sessionService.Login(ctx, &pb.UserLogin{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("Expected Unauthenticated error for wrong password, got: %v", err)
	}

	// Test Logout
	_, err = sessionService.Logout(ctx, &emptypb.Empty{})
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}
}

func TestListingService(t *testing.T) {
	store := storage.NewInMemoryStorage()
	service := NewListingService(store)
	ctx := context.Background()

	// Test CreateListing
	listing, err := service.CreateListing(ctx, &pb.ListingCreate{
		Title:       "iPhone 13",
		Description: "Great phone",
		Price:       999.99,
		Category:    "electronics",
		Condition:   "new",
		Location:    "New York",
	})
	if err != nil {
		t.Fatalf("CreateListing failed: %v", err)
	}
	if listing.Title != "iPhone 13" || listing.Price != 999.99 {
		t.Errorf("Listing data mismatch: got %+v", listing)
	}

	// Test GetListing
	retrievedListing, err := service.GetListing(ctx, &pb.GetListingRequest{Id: listing.Id})
	if err != nil {
		t.Fatalf("GetListing failed: %v", err)
	}
	if retrievedListing.Id != listing.Id {
		t.Errorf("Retrieved listing ID mismatch")
	}

	// Test GetListings with search
	listingsResp, err := service.GetListings(ctx, &pb.ListingsRequest{
		Search: "iPhone",
	})
	if err != nil {
		t.Fatalf("GetListings failed: %v", err)
	}
	if len(listingsResp.Listings) == 0 {
		t.Error("Expected to find listings with search term 'iPhone'")
	}

	// Test error cases
	_, err = service.GetListing(ctx, &pb.GetListingRequest{Id: 999})
	if status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound error for non-existent listing, got: %v", err)
	}
}

func TestOrderService(t *testing.T) {
	store := storage.NewInMemoryStorage()
	listingService := NewListingService(store)
	orderService := NewOrderService(store)
	ctx := context.Background()

	// Create a listing first
	listing, err := listingService.CreateListing(ctx, &pb.ListingCreate{
		Title:       "Test Product",
		Description: "Test description",
		Price:       100.0,
		Category:    "test",
		Condition:   "new",
	})
	if err != nil {
		t.Fatalf("Failed to create listing: %v", err)
	}

	// Test CreateOrder
	order, err := orderService.CreateOrder(ctx, &pb.OrderCreate{
		ListingId: listing.Id,
		Quantity:  2,
		ShippingAddress: &pb.Address{
			Street:  "123 Test St",
			City:    "Test City",
			Country: "USA",
		},
		BuyerNotes: "Test notes",
	})
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	if order.TotalPrice != 200.0 { // 100.0 * 2
		t.Errorf("Expected total price 200.0, got %f", order.TotalPrice)
	}

	// Test GetOrder
	retrievedOrder, err := orderService.GetOrder(ctx, &pb.GetOrderRequest{Id: order.Id})
	if err != nil {
		t.Fatalf("GetOrder failed: %v", err)
	}
	if retrievedOrder.Id != order.Id {
		t.Errorf("Retrieved order ID mismatch")
	}

	// Test UpdateOrderStatus
	updatedOrder, err := orderService.UpdateOrderStatus(ctx, &pb.UpdateOrderStatusRequest{
		Id:     order.Id,
		Status: "shipped",
	})
	if err != nil {
		t.Fatalf("UpdateOrderStatus failed: %v", err)
	}
	if updatedOrder.Status != "shipped" {
		t.Errorf("Expected status 'shipped', got '%s'", updatedOrder.Status)
	}

	// Test CancelOrder
	cancelResp, err := orderService.CancelOrder(ctx, &pb.CancelOrderRequest{
		Id:           order.Id,
		CancelReason: "Test cancellation",
	})
	if err != nil {
		t.Fatalf("CancelOrder failed: %v", err)
	}
	if cancelResp.Order.Status != "cancelled" {
		t.Errorf("Expected cancelled status, got '%s'", cancelResp.Order.Status)
	}

	// Test GetOrders
	ordersResp, err := orderService.GetOrders(ctx, &pb.OrdersRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("GetOrders failed: %v", err)
	}
	if len(ordersResp.Orders) == 0 {
		t.Error("Expected to find orders")
	}

	// Test error cases
	_, err = orderService.UpdateOrderStatus(ctx, &pb.UpdateOrderStatusRequest{
		Id:     order.Id,
		Status: "invalid_status",
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument error for invalid status, got: %v", err)
	}
}
