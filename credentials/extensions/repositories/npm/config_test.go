package npm_test

import (
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/npm"
	identity "github.com/mandelsoft/ctxmgmt/credentials/identity/npm"
	"github.com/mandelsoft/ctxmgmt/utils"
)

var _ = Describe("Config deserialization Test Environment", func() {
	It("read .npmrc", func() {
		ctx := credentials.New()
		repo := Must(npm.NewRepository(ctx, "testdata/.npmrc"))
		Expect(Must(repo.LookupCredentials("registry.npmjs.org")).Properties()).To(Equal(utils.Properties{identity.ATTR_TOKEN: "npm_TOKEN"}))
		Expect(Must(repo.LookupCredentials("npm.registry.acme.com/api/npm")).Properties()).To(Equal(utils.Properties{identity.ATTR_TOKEN: "bearer_TOKEN"}))
	})

	It("propagates credentials", func() {
		ctx := credentials.New()
		spec := npm.NewRepositorySpec("testdata/.npmrc")
		_ = Must(ctx.RepositoryForSpec(spec))
		id := Must(identity.GetConsumerId("registry.npmjs.org", "pkg"))
		creds := Must(credentials.CredentialsForConsumer(ctx, id))
		Expect(creds).NotTo(BeNil())
		Expect(creds.GetProperty(identity.ATTR_TOKEN)).To(Equal("npm_TOKEN"))
	})

	It("has description", func() {
		ctx := credentials.New()
		t := ctx.RepositoryTypes().GetType(npm.TypeV1)
		Expect(t).NotTo(BeNil())
		Expect(t.Description()).NotTo(Equal(""))
	})
})
