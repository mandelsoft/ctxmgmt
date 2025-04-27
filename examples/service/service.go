package service

import (
	"fmt"

	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/examples/service/identity"
)

type ServiceClient struct {
	address  string
	username string
	password string
	cert     string
}

// --- begin service ---
func NewServiceClient(ctx credentials.Context, addr string) (*ServiceClient, error) {
	creds, err := credentials.CredentialsForConsumer(ctx, identity.GetConsumerId(addr), identity.IdentityMatcher)
	if err != nil {
		return nil, err
	}

	s := &ServiceClient{
		address: addr,
	}

	if creds != nil {
		s.username = creds.GetProperty(credentials.ATTR_USERNAME)
		s.password = creds.GetProperty(credentials.ATTR_PASSWORD)
		s.cert = creds.GetProperty(credentials.ATTR_CERTIFICATE)
	}
	return s, nil
}

// --- end service ---

func (c *ServiceClient) Describe() {
	fmt.Printf("service client:\n")
	fmt.Printf(" username:    %s\n", c.username)
	fmt.Printf(" password:    %s\n", c.password)
	fmt.Printf(" certificate: %s\n", c.cert)
}
