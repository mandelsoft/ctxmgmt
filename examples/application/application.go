package application

import (
	"fmt"

	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/examples/helper"
	"github.com/mandelsoft/ctxmgmt/examples/myconfig"
	"github.com/mandelsoft/ctxmgmt/examples/service"
	"github.com/mandelsoft/goutils/errors"
)

// --- begin application ---
type Config struct {
	address string
}

func (a *Config) SetServiceAddress(addr string) {
	a.address = addr
}

type Application struct {
	Config
	client *service.ServiceClient
}

var _ myconfig.ConfigTarget = (*Application)(nil)

func NewApplication(ctx credentials.Context) (*Application, error) {
	app := &Application{}

	err := ctx.ConfigContext().ApplyAllTo(&app.Config)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot apply to application config")
	}
	app.client, err = service.NewServiceClient(ctx, app.address)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create service client")
	}
	return app, nil
}

// --- end application ---

func (a *Application) Describe() {
	helper.Output("app info", func() {
		fmt.Printf("service address: %s\n", a.address)
		a.client.Describe()
	})
}
