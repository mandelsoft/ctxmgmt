package internal

import (
	datalog "github.com/mandelsoft/datacontext/logging"
)

var Realm = datalog.DefineSubRealm("configuration management", "config")

var Logger = datalog.DynamicLogger(Realm)
