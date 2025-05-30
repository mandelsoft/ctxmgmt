package vault

import (
	"strings"

	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
	"github.com/mandelsoft/ctxmgmt/credentials/identity/vault"
	"github.com/mandelsoft/ctxmgmt/utils/listformat"
)

func init() {
	info := cpi.DefaultContext.ConsumerIdentityMatchers().GetInfo(vault.CONSUMER_TYPE)
	idx := strings.Index(info.Description, "\n")
	desc := `
This repository type can be used to access credentials stored in a HashiCorp
Vault. 

It provides access to list of secrets stored under a dedicated path in
a vault namespace. This list can either explicitly be specified, or
it is taken from the metadata of a specified secret.

The following custom metadata attributes are evaluated:
- <code>` + CUSTOM_SECRETS + `</code> this attribute may contain a comma separated list of
  vault secrets, which should be exposed by this repository instance.
  The names are evaluated under the path prefix used for the repository.
- <code>` + CUSTOM_CONSUMERID + `</code> this attribute may contain a JSON encoded
  consumer id , this secret should be assigned to.
- <code>type</code> if no special attribute is defined this attribute 
  indicated to use the complete custom metadata as consumer id.

It uses the ` + vault.CONSUMER_TYPE + ` identity matcher and consumer type
to requests credentials for the access.
` + info.Description[idx:] + `

It requires the following credential attributes:

` + info.CredentialAttributes

	usage = desc
}

var usage string

var format = `
The repository specification supports the following fields:
` + listformat.FormatListElements("", listformat.StringElementDescriptionList{
	"serverURL", "*string* (required): the URL of the vault instance",
	"namespace", "*string* (optional): the namespace used to evaluate secrets",
	"mountPath", "*string* (optional): the mount path to use (default: secrets)",
	"path", "*string* (optional): the path prefix used to lookup secrets",
	"secrets", "*[]string* (optional): list of secrets",
	"propagateConsumerIdentity", "*bool*(optional): evaluate metadata for consumer id propagation",
}) + `
If the secrets list is empty, all secret entries found in the given path
is read.
`
