package vault

import (
	"context"
	"sync"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	identity "github.com/mandelsoft/ctxmgmt/credentials/identity/vault"
	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/maputils"

	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
)

type AuthMethod interface {
	GetName() string
	Validate(creds cpi.Credentials) error
	GetToken(ctx context.Context, client *vault.Client, ns string, creds cpi.Credentials) (string, error)
}

type AuthMethods struct {
	lock    sync.Mutex
	methods map[string]AuthMethod
}

var methods = NewAuthMethods()

func NewAuthMethods() *AuthMethods {
	return &AuthMethods{
		methods: map[string]AuthMethod{},
	}
}

func (r *AuthMethods) Register(m AuthMethod) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.methods[m.GetName()] = m
}

func (r *AuthMethods) Get(name string) AuthMethod {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.methods[name]
}

func (r *AuthMethods) Names() []string {
	r.lock.Lock()
	defer r.lock.Unlock()

	return maputils.OrderedKeys(r.methods)
}

func RegisterAuthMethod(m AuthMethod) {
	methods.Register(m)
}

////////////////////////////////////////////////////////////////////////////////

func init() {
	RegisterAuthMethod(&approle{})
	RegisterAuthMethod(&token{})
}

////////////////////////////////////////////////////////////////////////////////

type approle struct{}

var _ AuthMethod = (*approle)(nil)

func (a *approle) GetName() string {
	return identity.AUTH_APPROLE
}

func (a *approle) Validate(creds cpi.Credentials) error {
	if !creds.ExistsProperty(identity.ATTR_ROLEID) {
		return errors.ErrRequired("credential property", identity.ATTR_ROLEID, a.GetName())
	}
	if !creds.ExistsProperty(identity.ATTR_SECRETID) {
		return errors.ErrRequired("credential property", identity.ATTR_SECRETID, a.GetName())
	}
	return nil
}

func (a *approle) GetToken(ctx context.Context, client *vault.Client, ns string, creds cpi.Credentials) (string, error) {
	req := schema.AppRoleLoginRequest{
		RoleId:   creds.GetProperty(identity.ATTR_ROLEID),
		SecretId: creds.GetProperty(identity.ATTR_SECRETID),
	}
	resp, err := client.Auth.AppRoleLogin(
		ctx,
		req,
		vault.WithNamespace(ns),
	)
	if err != nil {
		return "", err
	}
	return resp.Auth.ClientToken, nil
}

////////////////////////////////////////////////////////////////////////////////

type token struct{}

var _ AuthMethod = (*token)(nil)

func (a *token) GetName() string {
	return identity.AUTH_TOKEN
}

func (a *token) Validate(creds cpi.Credentials) error {
	if !creds.ExistsProperty(identity.ATTR_TOKEN) {
		return errors.ErrRequired("credential property", identity.ATTR_TOKEN, a.GetName())
	}
	return nil
}

func (a *token) GetToken(ctx context.Context, client *vault.Client, ns string, creds cpi.Credentials) (string, error) {
	return creds.GetProperty(identity.ATTR_TOKEN), nil
}
