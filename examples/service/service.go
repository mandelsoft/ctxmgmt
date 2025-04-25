package service

import (
	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/credentials/identity/hostpath"
	"github.com/mandelsoft/ctxmgmt/examples/service/identity"
)

const CONSUMER_ID = "service.example.acme.com"

type ServiceClient struct {
	address  string
	username string
	password string
}

func NewServiceClient(ctx credentials.Context, addr string) (*ServiceClient, error) {
	creds, err := credentials.CredentialsForConsumer(ctx, identity.GetConsumerId(addr), hostpath.IdentityMatcher(CONSUMER_ID))
	if err != nil {
		return nil, err
	}

	return &ServiceClient{
		address:  addr,
		username: creds.GetProperty(credentials.ATTR_USERNAME),
		password: creds.GetProperty(credentials.ATTR_PASSWORD),
	}, nil
}
