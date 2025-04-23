package vault

import (
	ctxlog "github.com/mandelsoft/datacontext/logging"
)

var (
	REALM = ctxlog.DefineSubRealm("HashiCorp Vault Access", "credentials", "vault")
	log   = ctxlog.DynamicLogger(REALM)
)
