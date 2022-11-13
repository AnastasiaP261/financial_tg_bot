package reports

import (
	"context"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/utils/logs"
	pkg "gitlab.ozon.dev/apetrichuk/financial-tg-bot/pkg/api/financial-tg-bot"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type configGetter interface {
	GRPCPortMessages() int
}

type Sender interface {
	SendReport(ctx context.Context, userID int64, text string, img []byte) error
}

type server struct {
	pkg.MessagesServiceServer
	sender Sender
}

func Register(conf configGetter, sender Sender) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPCPortMessages()))
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	logs.Info(
		"server listening",
		zap.Any("addr", lis.Addr()),
	)

	s := grpc.NewServer()

	pkg.RegisterMessagesServiceServer(s, &server{sender: sender})
	err = s.Serve(lis)
	if err != nil {
		return errors.Wrap(err, "Serve")
	}

	return nil
}
