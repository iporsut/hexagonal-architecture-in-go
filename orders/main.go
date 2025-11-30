package orders

import (
	"context"
)

type Product struct {
	ID    string
	Name  string
	Price float64
}

type Oder struct {
	ID       string
	UserID   string
	Items    []OrderItem
	TotalAmt float64
}

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

type OrderApp interface {
	ImportProducts(ctx context.Context, products []Product) error
	AddItemToCart(ctx context.Context, userID string, productID string, quantity int) error
	// NOTE: ingore inventory check for simplicity
	PlaceOrder(ctx context.Context, userID string) error
}

type OrderRepo interface {
	SaveOrder(ctx context.Context, order Oder) error
}

type CartRepo interface {
	AddItemToCart(ctx context.Context, userID string, item OrderItem) error
	GetCartItems(ctx context.Context, userID string) ([]OrderItem, error)
	ClearCart(ctx context.Context, userID string) error
}

type ProductRepo interface {
	ImportProducts(ctx context.Context, products []Product) error
	GetProductByID(ctx context.Context, productID string) (Product, error)
}

type OrderPlacedNotifier interface {
	NotifyOrderPlaced(ctx context.Context, order Oder) error
}

type orderApp struct {
	productRepo         ProductRepo
	cartRepo            CartRepo
	orderRepo           OrderRepo
	orderPlacedNotifier OrderPlacedNotifier
}

func (o *orderApp) ImportProducts(ctx context.Context, products []Product) error {
	return o.productRepo.ImportProducts(ctx, products)
}

func (o *orderApp) AddItemToCart(ctx context.Context, userID string, productID string, quantity int) error {
	product, err := o.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}
	item := OrderItem{
		ProductID: product.ID,
		Quantity:  quantity,
		Price:     product.Price,
	}
	err = o.cartRepo.AddItemToCart(ctx, userID, item)
	if err != nil {
		return err
	}
	return nil
}

func (o *orderApp) PlaceOrder(ctx context.Context, userID string) error {
	cartItems, err := o.cartRepo.GetCartItems(ctx, userID)
	if err != nil {
		return err
	}

	order := Oder{
		ID:     "order123", // In real app, generate unique ID
		UserID: userID,
		Items:  cartItems,
	}

	var total float64
	for _, item := range cartItems {
		total += item.Price * float64(item.Quantity)
	}
	order.TotalAmt = total

	err = o.orderRepo.SaveOrder(ctx, order)
	if err != nil {
		return err
	}

	err = o.cartRepo.ClearCart(ctx, userID)
	if err != nil {
		return err
	}

	if o.orderPlacedNotifier != nil {
		err = o.orderPlacedNotifier.NotifyOrderPlaced(ctx, order)
		if err != nil {
			return err
		}
	}

	return nil
}

type OrderAppOption func(*orderApp)

func WithProductRepo(repo ProductRepo) OrderAppOption {
	return func(o *orderApp) {
		o.productRepo = repo
	}
}

func WithCartRepo(repo CartRepo) OrderAppOption {
	return func(o *orderApp) {
		o.cartRepo = repo
	}
}

func WithOrderRepo(repo OrderRepo) OrderAppOption {
	return func(o *orderApp) {
		o.orderRepo = repo
	}
}

func WithOrderPlacedNotifier(notifier OrderPlacedNotifier) OrderAppOption {
	return func(o *orderApp) {
		o.orderPlacedNotifier = notifier
	}
}

func NewOrderApp(opts ...OrderAppOption) *orderApp {
	o := &orderApp{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Notification adapter for OrderPlacedNotifier can be added here if needed

type BroadcastOrderPlacedNotifier struct {
	notifiers []OrderPlacedNotifier
}

func NewBroadcastOrderPlacedNotifier(notifiers []OrderPlacedNotifier) *BroadcastOrderPlacedNotifier {
	return &BroadcastOrderPlacedNotifier{
		notifiers: notifiers,
	}
}

func (b *BroadcastOrderPlacedNotifier) NotifyOrderPlaced(ctx context.Context, order Oder) error {
	for _, notifier := range b.notifiers {
		err := notifier.NotifyOrderPlaced(ctx, order)
		if err != nil {
			return err
		}
	}
	return nil
}

type SMSNotifier struct {
}

func (s *SMSNotifier) NotifyOrderPlaced(ctx context.Context, order Oder) error {
	// Simulate sending SMS
	return nil
}

type EmailNotifier struct {
}

func (e *EmailNotifier) NotifyOrderPlaced(ctx context.Context, order Oder) error {
	// Simulate sending Email
	return nil
}

func main() {
	notifier := NewBroadcastOrderPlacedNotifier([]OrderPlacedNotifier{
		&SMSNotifier{},
		&EmailNotifier{},
	})

	orderApp := NewOrderApp(
		WithOrderPlacedNotifier(notifier),
	)
}
