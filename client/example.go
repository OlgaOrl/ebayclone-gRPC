package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "ebayclone-grpc/proto"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create service clients
	userClient := pb.NewUserServiceClient(conn)
	sessionClient := pb.NewSessionServiceClient(conn)
	listingClient := pb.NewListingServiceClient(conn)
	orderClient := pb.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Println("=== eBayClone gRPC Client Example ===")

	// 1. Create a user
	log.Println("\n1. Creating user...")
	user, err := userClient.CreateUser(ctx, &pb.UserCreate{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return
	}
	log.Printf("Created user: ID=%d, Username=%s, Email=%s", user.Id, user.Username, user.Email)

	// 2. Login
	log.Println("\n2. Logging in...")
	loginResp, err := sessionClient.Login(ctx, &pb.UserLogin{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Printf("Failed to login: %v", err)
		return
	}
	log.Printf("Login successful, token: %s", loginResp.Token[:20]+"...")

	// 3. Get user
	log.Println("\n3. Getting user...")
	retrievedUser, err := userClient.GetUser(ctx, &pb.GetUserRequest{Id: user.Id})
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return
	}
	log.Printf("Retrieved user: ID=%d, Username=%s", retrievedUser.Id, retrievedUser.Username)

	// 4. Create a listing
	log.Println("\n4. Creating listing...")
	listing, err := listingClient.CreateListing(ctx, &pb.ListingCreate{
		Title:       "iPhone 13 Pro Max",
		Description: "Brand new, still in box",
		Price:       999.99,
		Category:    "electronics",
		Condition:   "new",
		Location:    "New York, NY",
	})
	if err != nil {
		log.Printf("Failed to create listing: %v", err)
		return
	}
	log.Printf("Created listing: ID=%d, Title=%s, Price=$%.2f", listing.Id, listing.Title, listing.Price)

	// 5. Get listings
	log.Println("\n5. Getting listings...")
	listingsResp, err := listingClient.GetListings(ctx, &pb.ListingsRequest{
		Search:   "iPhone",
		PriceMin: 500,
		PriceMax: 1500,
	})
	if err != nil {
		log.Printf("Failed to get listings: %v", err)
		return
	}
	log.Printf("Found %d listings", len(listingsResp.Listings))

	// 6. Create an order
	log.Println("\n6. Creating order...")
	order, err := orderClient.CreateOrder(ctx, &pb.OrderCreate{
		ListingId: listing.Id,
		Quantity:  1,
		ShippingAddress: &pb.Address{
			Street:  "123 Main St",
			City:    "New York",
			State:   "NY",
			ZipCode: "10001",
			Country: "USA",
		},
		BuyerNotes: "Please deliver after 5 PM",
	})
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		return
	}
	log.Printf("Created order: ID=%d, Status=%s, Total=$%.2f", order.Id, order.Status, order.TotalPrice)

	// 7. Get orders
	log.Println("\n7. Getting orders...")
	ordersResp, err := orderClient.GetOrders(ctx, &pb.OrdersRequest{
		Page:  1,
		Limit: 10,
	})
	if err != nil {
		log.Printf("Failed to get orders: %v", err)
		return
	}
	log.Printf("Found %d orders (Total: %d, Pages: %d)", len(ordersResp.Orders), ordersResp.Pagination.Total, ordersResp.Pagination.Pages)

	// 8. Update order status
	log.Println("\n8. Updating order status...")
	updatedOrder, err := orderClient.UpdateOrderStatus(ctx, &pb.UpdateOrderStatusRequest{
		Id:     order.Id,
		Status: "confirmed",
	})
	if err != nil {
		log.Printf("Failed to update order status: %v", err)
		return
	}
	log.Printf("Updated order status to: %s", updatedOrder.Status)

	// 9. Cancel order
	log.Println("\n9. Cancelling order...")
	cancelResp, err := orderClient.CancelOrder(ctx, &pb.CancelOrderRequest{
		Id:           order.Id,
		CancelReason: "Changed my mind",
	})
	if err != nil {
		log.Printf("Failed to cancel order: %v", err)
		return
	}
	log.Printf("Order cancelled: %s", cancelResp.Message)

	// 10. Logout
	log.Println("\n10. Logging out...")
	_, err = sessionClient.Logout(ctx, &emptypb.Empty{})
	if err != nil {
		log.Printf("Failed to logout: %v", err)
		return
	}
	log.Println("Logout successful")

	log.Println("\n=== All operations completed successfully! ===")
}
