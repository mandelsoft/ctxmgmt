package internal

import (
	datalog "github.com/mandelsoft/ctxmgmt/logging"
)

var Realm = datalog.DefineSubRealm("configuration management", "config")

var Logger = datalog.DynamicLogger(Realm)
