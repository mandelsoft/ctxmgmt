package application

import (
	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/examples/myconfig"
	"github.com/mandelsoft/ctxmgmt/examples/service"
)

type Config struct {
	address string
}

type Application struct {
	address string
	client  *service.ServiceClient
}

var _ myconfig.ConfigTarget = (*Application)(nil)

func NewApplication(ctx credentials.Context) (*Application, error) {
	cfg := &Config{}

	ctx.ConfigContext().ApplyTo(0, cfg)
	c, err := service.NewServiceClient(ctx, cfg.address)
	if err != nil {
		return nil, err
	}
	return &Application{cfg.address, c}, nil
}

func (a *Application) SetServiceAddress(addr string) {
	a.address = addr
}
