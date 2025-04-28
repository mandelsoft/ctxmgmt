package main

import (
	"fmt"

	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/credentials/identity/hostpath"
	"github.com/mandelsoft/ctxmgmt/examples/helper"
	"github.com/mandelsoft/goutils/errors"
)

func Configure(ctx credentials.Context) {
	// the most simple usecase is to directly define credentials
	// for a dedicated usage scenario at the credential context.
	// A consumption scenario is defined by a ConsumerIdentity.
	// It always describes the scenario type (demo) and additional
	// attributes which concretize the scenario.
	// Credentials consist of a set of credential attributes,
	// like username or password.

	// --- begin setting credentials ---
	defid := credentials.NewConsumerIdentity("demo", hostpath.ID_HOSTNAME, "localhost", hostpath.ID_PATHPREFIX, "repositories")

	creds := credentials.CredentialsFromList(credentials.ATTR_USERNAME, "testuser", credentials.ATTR_PASSWORD, "testpass")
	// creds := crednetials.SimpleCredentials("testuser", "testpass")

	ctx.SetCredentialsForConsumer(defid, creds)
	// --- end setting credentials ---

	helper.Output("definition id", func() {
		fmt.Printf("setting credentials for %s\n", defid)
	})
}

func BasicCredentialManagement() error {
	// credentials are handled by the credential context.
	// --- begin default context ---
	ctx := credentials.DefaultContext()
	// --- end default context ---

	Configure(ctx)

	// now, we can try to request credentials.
	// First, we describe our intended scenario.
	// We want to access some element of this demo scenario under the path
	// repositories/first

	// --- begin request id ---
	cid := credentials.NewConsumerIdentity("demo", hostpath.ID_HOSTNAME, "localhost", hostpath.ID_PATHPREFIX, "repositories/first")

	helper.Output("request id", func() {
		fmt.Printf("requesting credentials for %s\n", cid)
	})
	// --- end request id ---

	// with this consumer id we can query credentials from the context.
	// because we use an ad-hoc scenario type (demo), we explicitly specify an appropriate
	// identity matcher. We use the default hostpath matcher able to match path prefixes.

	// --- begin request credentials ---
	credsrc, err := ctx.GetCredentialsForConsumer(cid, hostpath.IdentityMatcher("demo"))
	if err != nil {
		return errors.Wrapf(err, "error matching consumer id %s", cid)
	}
	if credsrc == nil {
		return fmt.Errorf("no crendials found for %s", cid)
	}
	// --- end request credentials ---

	// the provided result represents a dynamic credential representation. It might change
	// if the underlying credential source changes the credentials.
	// Is is used to retrieve concrete credential attributes.

	// --- begin concrete values ---
	creds, err := credsrc.Credentials(ctx)
	if err != nil {
		return errors.Wrapf(err, "error getting credentials for %s", cid)
	}
	// --- end concrete values ---

	// the provided Credentials object can be used to query the provided
	// credential attributes.
	helper.Output("credentials", func() {
		fmt.Printf("found crednetials for %s\n", creds)
	})

	// If we would use a path not matching the originally configured prefix
	// no credentials are provided.

	// --- begin changed ---
	cid[hostpath.ID_PATHPREFIX] = "otherpath"

	creds, err = credentials.CredentialsForConsumer(ctx, cid, hostpath.IdentityMatcher("demo"))
	if creds != nil {
		return fmt.Errorf("Ooops, found credentials for %s", cid)
	}
	// --- end changed ---

	return nil
}
