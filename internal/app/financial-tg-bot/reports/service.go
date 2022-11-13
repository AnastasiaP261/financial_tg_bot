package reports

import (
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/apetrichuk/financial-tg-bot/internal/logs"
	pkg "gitlab.ozon.dev/apetrichuk/financial-tg-bot/pkg/api/financial-tg-bot"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	timestampFormat = time.StampNano
)

type configGetter interface {
	GRPCPortMessages() int
}

type server struct {
	pkg.MessagesServiceServer
}

func Register(conf configGetter) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GRPCPortMessages()))
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	logs.Info(
		"server listening",
		zap.Any("addr", lis.Addr()),
	)

	s := grpc.NewServer()

	pkg.RegisterMessagesServiceServer(s, &server{})
	err = s.Serve(lis)
	if err != nil {
		return errors.Wrap(err, "Serve")
	}

	return nil
}
