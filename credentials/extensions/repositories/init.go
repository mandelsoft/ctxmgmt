package repositories

import (
	_ "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/aliases"
	_ "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/directcreds"
	_ "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/dockerconfig"
	_ "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/memory"
	_ "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/memory/config"
	_ "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/npm"
	_ "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/vault"
)
