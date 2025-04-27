package main

import (
	"fmt"
	"os"

	"github.com/mandelsoft/ctxmgmt/config/extensions/config"
	credcfg "github.com/mandelsoft/ctxmgmt/credentials/config"
	"github.com/mandelsoft/ctxmgmt/examples/myconfig"
	"github.com/mandelsoft/ctxmgmt/examples/service/identity"
	"github.com/mandelsoft/goutils/generics"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config config type":
			fmt.Printf("%s", config.ConfigTypeV1)
		case "credentials config type":
			fmt.Printf("%s", credcfg.ConfigTypeV1)
		case "my config type":
			fmt.Printf("%s", myconfig.TYPE)
		case "consumer type":
			fmt.Printf("%s", identity.CONSUMER_TYPE)
		case "config object struct":
			fmt.Printf("%s", generics.TypeOf[myconfig.Config]().Name())
		case "config target":
			fmt.Printf("%s", generics.TypeOf[myconfig.ConfigTarget]().Name())
		default:
			panic(fmt.Sprintf("unkonwn key %q", os.Args[1]))
		}
	}
}
