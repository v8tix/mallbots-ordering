package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	"github.com/v8tix/mallbots-ordering-proto/pb"
	"github.com/v8tix/mallbots-ordering/internal/application"
	"github.com/v8tix/mallbots-ordering/internal/application/commands"
	"github.com/v8tix/mallbots-ordering/internal/application/queries"
	"github.com/v8tix/mallbots-ordering/internal/domain"
)

type server struct {
	app application.App
	pb.UnimplementedOrderingServiceServer
}

var _ pb.OrderingServiceServer = (*server)(nil)

func RegisterServer(app application.App, registrar grpc.ServiceRegistrar) error {
	pb.RegisterOrderingServiceServer(registrar, server{app: app})
	return nil
}

func (s server) CreateOrder(ctx context.Context, request *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	id := uuid.New().String()

	items := make([]domain.Item, len(request.Items))
	for i, item := range request.Items {
		items[i] = s.itemToDomain(item)
	}

	err := s.app.CreateOrder(ctx, commands.CreateOrder{
		ID:         id,
		CustomerID: request.GetCustomerId(),
		PaymentID:  request.GetPaymentId(),
		Items:      items,
	})

	return &pb.CreateOrderResponse{Id: id}, err
}

func (s server) CancelOrder(ctx context.Context, request *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	err := s.app.CancelOrder(ctx, commands.CancelOrder{ID: request.GetId()})

	return &pb.CancelOrderResponse{}, err
}

func (s server) ReadyOrder(ctx context.Context, request *pb.ReadyOrderRequest) (*pb.ReadyOrderResponse, error) {
	err := s.app.ReadyOrder(ctx, commands.ReadyOrder{ID: request.GetId()})
	return &pb.ReadyOrderResponse{}, err
}

func (s server) CompleteOrder(ctx context.Context, request *pb.CompleteOrderRequest) (*pb.CompleteOrderResponse, error) {
	err := s.app.CompleteOrder(ctx, commands.CompleteOrder{ID: request.GetId()})
	return &pb.CompleteOrderResponse{}, err
}

func (s server) GetOrder(ctx context.Context, request *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.app.GetOrder(ctx, queries.GetOrder{ID: request.GetId()})
	if err != nil {
		return nil, err
	}

	return &pb.GetOrderResponse{
		Order: s.orderFromDomain(order),
	}, nil
}

func (s server) orderFromDomain(order *domain.Order) *pb.Order {
	items := make([]*pb.OrderingItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = s.itemFromDomain(item)
	}

	return &pb.Order{
		Id:         order.ID(),
		CustomerId: order.CustomerID,
		PaymentId:  order.PaymentID,
		Items:      items,
		Status:     order.Status.String(),
	}
}

func (s server) itemToDomain(item *pb.OrderingItem) domain.Item {
	return domain.Item{
		ProductID:   item.GetProductId(),
		StoreID:     item.GetStoreId(),
		StoreName:   item.GetStoreName(),
		ProductName: item.GetProductName(),
		Price:       item.GetPrice(),
		Quantity:    int(item.GetQuantity()),
	}
}

func (s server) itemFromDomain(item domain.Item) *pb.OrderingItem {
	return &pb.OrderingItem{
		StoreId:     item.StoreID,
		ProductId:   item.ProductID,
		StoreName:   item.StoreName,
		ProductName: item.ProductName,
		Price:       item.Price,
		Quantity:    int32(item.Quantity),
	}
}
