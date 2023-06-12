package grpc

import (
	"context"
	"database/sql"

	"google.golang.org/grpc"

	"github.com/v8tix/eda/di"
	"github.com/v8tix/mallbots-ordering-proto/pb"
	"github.com/v8tix/mallbots-ordering/internal/application"
)

type serverTx struct {
	c di.Container
	pb.UnimplementedOrderingServiceServer
}

var _ pb.OrderingServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	pb.RegisterOrderingServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateOrder(ctx context.Context, request *pb.CreateOrderRequest) (resp *pb.CreateOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.CreateOrder(ctx, request)
}

func (s serverTx) GetOrder(ctx context.Context, request *pb.GetOrderRequest) (resp *pb.GetOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.GetOrder(ctx, request)
}

func (s serverTx) CancelOrder(ctx context.Context, request *pb.CancelOrderRequest) (resp *pb.CancelOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.CancelOrder(ctx, request)
}

func (s serverTx) ReadyOrder(ctx context.Context, request *pb.ReadyOrderRequest) (resp *pb.ReadyOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.ReadyOrder(ctx, request)
}

func (s serverTx) CompleteOrder(ctx context.Context, request *pb.CompleteOrderRequest) (resp *pb.CompleteOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, "tx").(*sql.Tx))

	next := server{app: di.Get(ctx, "app").(application.App)}

	return next.CompleteOrder(ctx, request)
}

func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
