package directcreds

import (
	"github.com/mandelsoft/ctxmgmt/utils/listformat"
)

var usage = `
This repository type can be used to specify a single inline credential
set. The default name is the empty string or <code>` + Type + `</code>.`

var format = `The repository specification supports the following fields:
` + listformat.FormatListElements("", listformat.StringElementDescriptionList{
	"properties", "*map[string]string*: direct credential fields",
})
