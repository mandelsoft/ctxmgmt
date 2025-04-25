package ctxmgmt

import (
	"github.com/mandelsoft/ctxmgmt/logging"
)

var Realm = logging.DefineSubRealm("context lifecycle", "context")

var Logger = logging.DynamicLogger(Realm)

func Debug(c Context, msg string, keypairs ...interface{}) {
	c.LoggingContext().Logger(Realm).Debug(msg, append(keypairs, "id", c.GetId())...)
}
