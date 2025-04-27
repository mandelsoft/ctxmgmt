package main

import (
	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/config/cfgutils"
	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/examples/application"
	"github.com/mandelsoft/goutils/errors"
)

// --- begin main ---
func RunApplication() error {

	ctx := credentials.New(ctxmgmt.MODE_SHARED)

	err := cfgutils.Configure(ctx, ".appconfig")
	if err != nil {
		return errors.Wrap(err, "reading configuration")
	}

	app, err := application.NewApplication(ctx)
	if err != nil {
		return errors.Wrap(err, "error creating application")
	}

	app.Describe()
	return nil
}

// --- end main ---
