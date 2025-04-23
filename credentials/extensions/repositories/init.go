package repositories

import (
	_ "github.com/mandelsoft/datacontext/credentials/extensions/repositories/aliases"
	_ "github.com/mandelsoft/datacontext/credentials/extensions/repositories/directcreds"
	_ "github.com/mandelsoft/datacontext/credentials/extensions/repositories/dockerconfig"
	_ "github.com/mandelsoft/datacontext/credentials/extensions/repositories/memory"
	_ "github.com/mandelsoft/datacontext/credentials/extensions/repositories/memory/config"
	_ "github.com/mandelsoft/datacontext/credentials/extensions/repositories/npm"
	_ "github.com/mandelsoft/datacontext/credentials/extensions/repositories/vault"
)
