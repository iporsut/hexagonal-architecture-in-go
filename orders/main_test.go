package orders

import (
	"context"
	"testing"
)

type InMemoryProductRepo struct {
	products map[string]Product
}

func NewInMemoryProductRepo() *InMemoryProductRepo {
	return &InMemoryProductRepo{
		products: make(map[string]Product),
	}
}

func (r *InMemoryProductRepo) ImportProducts(ctx context.Context, products []Product) error {
	for _, p := range products {
		r.products[p.ID] = p
	}
	return nil
}

func (r *InMemoryProductRepo) GetProductByID(ctx context.Context, productID string) (Product, error) {
	if p, exists := r.products[productID]; exists {
		return p, nil
	}
	return Product{}, nil
}

type InMemoryCartRepo struct {
	carts map[string][]OrderItem
}

func NewInMemoryCartRepo() *InMemoryCartRepo {
	return &InMemoryCartRepo{
		carts: make(map[string][]OrderItem),
	}
}

func (r *InMemoryCartRepo) AddItemToCart(ctx context.Context, userID string, item OrderItem) error {
	r.carts[userID] = append(r.carts[userID], item)
	return nil
}

func (r *InMemoryCartRepo) GetCartItems(ctx context.Context, userID string) ([]OrderItem, error) {
	return r.carts[userID], nil
}

func (r *InMemoryCartRepo) ClearCart(ctx context.Context, userID string) error {
	delete(r.carts, userID)
	return nil
}

type InMemoryOrderRepo struct {
	orders []Oder
}

func NewInMemoryOrderRepo() *InMemoryOrderRepo {
	return &InMemoryOrderRepo{
		orders: []Oder{},
	}
}

func (r *InMemoryOrderRepo) SaveOrder(ctx context.Context, order Oder) error {
	r.orders = append(r.orders, order)
	return nil
}

type MockOrderPlacedNotifier struct {
	notifiedOrders []Oder
}

func NewMockOrderPlacedNotifier() *MockOrderPlacedNotifier {
	return &MockOrderPlacedNotifier{
		notifiedOrders: []Oder{},
	}
}

func (n *MockOrderPlacedNotifier) NotifyOrderPlaced(ctx context.Context, order Oder) error {
	n.notifiedOrders = append(n.notifiedOrders, order)
	return nil
}

func TestImportProducts(t *testing.T) {
	productRepo := NewInMemoryProductRepo()
	orderApp := NewOrderApp(
		WithProductRepo(productRepo),
	)

	products := []Product{
		{ID: "prod1", Name: "Product 1", Price: 10.0},
		{ID: "prod2", Name: "Product 2", Price: 20.0},
	}

	err := orderApp.ImportProducts(context.Background(), products)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(productRepo.products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(productRepo.products))
	}
}

func TestAddItemToCart(t *testing.T) {
	productRepo := NewInMemoryProductRepo()
	cartRepo := NewInMemoryCartRepo()
	orderApp := NewOrderApp(
		WithProductRepo(productRepo),
		WithCartRepo(cartRepo),
	)
	products := []Product{
		{ID: "prod1", Name: "Product 1", Price: 10.0},
	}
	err := orderApp.ImportProducts(context.Background(), products)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = orderApp.AddItemToCart(context.Background(), "user1", "prod1", 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(cartRepo.carts["user1"]) != 1 {
		t.Fatalf("expected 1 item in cart, got %d", len(cartRepo.carts["user1"]))
	}

	if cartRepo.carts["user1"][0].ProductID != "prod1" || cartRepo.carts["user1"][0].Quantity != 2 {
		t.Fatalf("cart item does not match expected values")
	}
}

func TestPlaceOrder(t *testing.T) {
	productRepo := NewInMemoryProductRepo()
	cartRepo := NewInMemoryCartRepo()
	orderRepo := NewInMemoryOrderRepo()
	notifier := NewMockOrderPlacedNotifier()
	orderApp := NewOrderApp(
		WithProductRepo(productRepo),
		WithCartRepo(cartRepo),
		WithOrderRepo(orderRepo),
		WithOrderPlacedNotifier(notifier),
	)

	products := []Product{
		{ID: "prod1", Name: "Product 1", Price: 10.0},
	}
	err := orderApp.ImportProducts(context.Background(), products)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = orderApp.AddItemToCart(context.Background(), "user1", "prod1", 2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = orderApp.PlaceOrder(context.Background(), "user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cartItems, err := orderApp.cartRepo.GetCartItems(context.Background(), "user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cartItems) != 0 {
		t.Fatalf("expected cart to be cleared, got %d items", len(cartItems))
	}

	if len(orderRepo.orders) != 1 {
		t.Fatalf("expected 1 order saved, got %d", len(orderRepo.orders))
	}

	if len(notifier.notifiedOrders) != 1 {
		t.Fatalf("expected 1 notification sent, got %d", len(notifier.notifiedOrders))
	}
}
