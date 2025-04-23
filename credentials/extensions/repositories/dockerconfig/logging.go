package dockerconfig

import (
	ctxlog "github.com/mandelsoft/datacontext/logging"
)

var REALM = ctxlog.DefineSubRealm("docker config handling as credential repository", "credentials/dockerconfig")
