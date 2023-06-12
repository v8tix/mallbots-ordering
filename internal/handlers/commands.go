package handlers

import (
	"context"

	"github.com/v8tix/eda/am"
	"github.com/v8tix/eda/ddd"
	"github.com/v8tix/mallbots-ordering-proto/pb"
	"github.com/v8tix/mallbots-ordering/internal/application"
	"github.com/v8tix/mallbots-ordering/internal/application/commands"
)

type commandHandlers struct {
	app application.App
}

func NewCommandHandlers(app application.App) ddd.CommandHandler[ddd.Command] {
	return commandHandlers{
		app: app,
	}
}

func RegisterCommandHandlers(subscriber am.RawMessageSubscriber, handlers am.RawMessageHandler) error {
	return subscriber.Subscribe(pb.CommandChannel, handlers, am.MessageFilter{
		pb.RejectOrderCommand,
		pb.ApproveOrderCommand,
	}, am.GroupName("ordering-commands"))
}

func (h commandHandlers) HandleCommand(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	switch cmd.CommandName() {
	case pb.RejectOrderCommand:
		return h.doRejectOrder(ctx, cmd)
	case pb.ApproveOrderCommand:
		return h.doApproveOrder(ctx, cmd)
	}

	return nil, nil
}

func (h commandHandlers) doRejectOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*pb.RejectOrder)

	return nil, h.app.RejectOrder(ctx, commands.RejectOrder{ID: payload.GetId()})
}

func (h commandHandlers) doApproveOrder(ctx context.Context, cmd ddd.Command) (ddd.Reply, error) {
	payload := cmd.Payload().(*pb.ApproveOrder)

	return nil, h.app.ApproveOrder(ctx, commands.ApproveOrder{
		ID:         payload.GetId(),
		ShoppingID: payload.GetShoppingId(),
	})
}
