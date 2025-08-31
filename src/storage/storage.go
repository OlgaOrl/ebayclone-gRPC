package storage

import (
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	pb "ebayclone-grpc/proto"
)

type Storage interface {
	// Users
	CreateUser(user *pb.User) error
	GetUser(id int32) (*pb.User, error)
	GetUserByEmail(email string) (*pb.User, error)
	UpdateUser(id int32, user *pb.User) error
	DeleteUser(id int32) error

	// Listings
	CreateListing(listing *pb.Listing) error
	GetListing(id int32) (*pb.Listing, error)
	GetListings(search string, priceMin, priceMax float64) ([]*pb.Listing, error)
	UpdateListing(id int32, listing *pb.Listing) error
	DeleteListing(id int32) error

	// Orders
	CreateOrder(order *pb.Order) error
	GetOrder(id int32) (*pb.Order, error)
	GetOrders(userID int32, status string, page, limit int32) ([]*pb.Order, int32, error)
	UpdateOrder(id int32, order *pb.Order) error
	DeleteOrder(id int32) error
}

type InMemoryStorage struct {
	mu       sync.RWMutex
	users    map[int32]*pb.User
	listings map[int32]*pb.Listing
	orders   map[int32]*pb.Order
	userID   int32
	listingID int32
	orderID  int32
	passwords map[int32]string // Store passwords separately for security
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users:     make(map[int32]*pb.User),
		listings:  make(map[int32]*pb.Listing),
		orders:    make(map[int32]*pb.Order),
		passwords: make(map[int32]string),
		userID:    1,
		listingID: 1,
		orderID:   1,
	}
}

// User methods
func (s *InMemoryStorage) CreateUser(user *pb.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if email already exists
	for _, existingUser := range s.users {
		if existingUser.Email == user.Email {
			return &UserExistsError{Email: user.Email}
		}
	}

	user.Id = s.userID
	s.users[s.userID] = user
	s.userID++
	return nil
}

func (s *InMemoryStorage) GetUser(id int32) (*pb.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, &NotFoundError{Resource: "User", ID: id}
	}
	return user, nil
}

func (s *InMemoryStorage) GetUserByEmail(email string) (*pb.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, &NotFoundError{Resource: "User", Email: email}
}

func (s *InMemoryStorage) UpdateUser(id int32, user *pb.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return &NotFoundError{Resource: "User", ID: id}
	}

	user.Id = id
	s.users[id] = user
	return nil
}

func (s *InMemoryStorage) DeleteUser(id int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return &NotFoundError{Resource: "User", ID: id}
	}

	delete(s.users, id)
	delete(s.passwords, id)
	return nil
}

func (s *InMemoryStorage) SetUserPassword(userID int32, password string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.passwords[userID] = password
}

func (s *InMemoryStorage) GetUserPassword(userID int32) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	password, exists := s.passwords[userID]
	return password, exists
}

// Custom error types
type NotFoundError struct {
	Resource string
	ID       int32
	Email    string
}

func (e *NotFoundError) Error() string {
	if e.Email != "" {
		return e.Resource + " with email " + e.Email + " not found"
	}
	return e.Resource + " not found"
}

type UserExistsError struct {
	Email string
}

func (e *UserExistsError) Error() string {
	return "User with email " + e.Email + " already exists"
}

// Listing methods
func (s *InMemoryStorage) CreateListing(listing *pb.Listing) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	listing.Id = s.listingID
	now := time.Now()
	listing.CreatedAt = timestamppb.New(now)
	listing.UpdatedAt = timestamppb.New(now)
	s.listings[s.listingID] = listing
	s.listingID++
	return nil
}

func (s *InMemoryStorage) GetListing(id int32) (*pb.Listing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	listing, exists := s.listings[id]
	if !exists {
		return nil, &NotFoundError{Resource: "Listing", ID: id}
	}
	return listing, nil
}

func (s *InMemoryStorage) GetListings(search string, priceMin, priceMax float64) ([]*pb.Listing, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*pb.Listing
	for _, listing := range s.listings {
		// Apply search filter
		if search != "" {
			if !contains(listing.Title, search) && !contains(listing.Description, search) {
				continue
			}
		}

		// Apply price filters
		if priceMin > 0 && listing.Price < priceMin {
			continue
		}
		if priceMax > 0 && listing.Price > priceMax {
			continue
		}

		result = append(result, listing)
	}
	return result, nil
}

func (s *InMemoryStorage) UpdateListing(id int32, listing *pb.Listing) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.listings[id]
	if !exists {
		return &NotFoundError{Resource: "Listing", ID: id}
	}

	listing.Id = id
	listing.CreatedAt = existing.CreatedAt
	listing.UpdatedAt = timestamppb.New(time.Now())
	s.listings[id] = listing
	return nil
}

func (s *InMemoryStorage) DeleteListing(id int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.listings[id]; !exists {
		return &NotFoundError{Resource: "Listing", ID: id}
	}

	delete(s.listings, id)
	return nil
}

// Order methods
func (s *InMemoryStorage) CreateOrder(order *pb.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	order.Id = s.orderID
	now := time.Now()
	order.CreatedAt = timestamppb.New(now)
	order.UpdatedAt = timestamppb.New(now)
	order.Status = "pending"
	s.orders[s.orderID] = order
	s.orderID++
	return nil
}

func (s *InMemoryStorage) GetOrder(id int32) (*pb.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, exists := s.orders[id]
	if !exists {
		return nil, &NotFoundError{Resource: "Order", ID: id}
	}
	return order, nil
}

func (s *InMemoryStorage) GetOrders(userID int32, status string, page, limit int32) ([]*pb.Order, int32, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var filtered []*pb.Order
	for _, order := range s.orders {
		// Apply filters
		if userID > 0 && order.UserId != userID {
			continue
		}
		if status != "" && order.Status != status {
			continue
		}
		filtered = append(filtered, order)
	}

	total := int32(len(filtered))

	// Apply pagination
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= total {
		return []*pb.Order{}, total, nil
	}
	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

func (s *InMemoryStorage) UpdateOrder(id int32, order *pb.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.orders[id]
	if !exists {
		return &NotFoundError{Resource: "Order", ID: id}
	}

	order.Id = id
	order.CreatedAt = existing.CreatedAt
	order.UpdatedAt = timestamppb.New(time.Now())
	s.orders[id] = order
	return nil
}

func (s *InMemoryStorage) DeleteOrder(id int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[id]; !exists {
		return &NotFoundError{Resource: "Order", ID: id}
	}

	delete(s.orders, id)
	return nil
}

// Helper function
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
