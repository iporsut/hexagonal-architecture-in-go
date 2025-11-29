package orders

import "context"

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
	return nil
}

func (o *orderApp) AddItemToCart(ctx context.Context, userID string, productID string, quantity int) error {
	return nil
}

func (o *orderApp) PlaceOrder(ctx context.Context, userID string) error {
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
