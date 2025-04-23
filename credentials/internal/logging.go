package internal

import (
	ctxlog "github.com/mandelsoft/datacontext/logging"
)

var (
	REALM = ctxlog.DefineSubRealm("Credentials", "credentials")
	log   = ctxlog.DynamicLogger(REALM)
)
