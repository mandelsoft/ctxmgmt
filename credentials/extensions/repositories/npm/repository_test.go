package npm_test

import (
	"encoding/json"
	"reflect"

	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/datacontext/credentials"
	"github.com/mandelsoft/datacontext/credentials/cpi"
	local "github.com/mandelsoft/datacontext/credentials/extensions/repositories/npm"
	identity "github.com/mandelsoft/datacontext/credentials/identity/npm"
	"github.com/mandelsoft/datacontext/utils"
	"github.com/mandelsoft/datacontext/utils/runtimefinalizer"
)

var _ = Describe("NPM config - .npmrc", func() {
	props := utils.Properties{
		identity.ATTR_TOKEN: "npm_TOKEN",
	}

	props2 := utils.Properties{
		identity.ATTR_TOKEN: "bearer_TOKEN",
	}

	var DefaultContext credentials.Context

	BeforeEach(func() {
		DefaultContext = credentials.New()
	})

	specdata := "{\"type\":\"NPMConfig\",\"npmrcFile\":\"testdata/.npmrc\"}"

	It("serializes repo spec", func() {
		spec := local.NewRepositorySpec("testdata/.npmrc")
		data := Must(json.Marshal(spec))
		Expect(data).To(Equal([]byte(specdata)))
	})

	It("deserializes repo spec", func() {
		spec := Must(DefaultContext.RepositorySpecForConfig([]byte(specdata), nil))
		Expect(reflect.TypeOf(spec).String()).To(Equal("*npm.RepositorySpec"))
		Expect(spec.(*local.RepositorySpec).NpmrcFile).To(Equal("testdata/.npmrc"))
	})

	It("resolves repository", func() {
		repo := Must(DefaultContext.RepositoryForConfig([]byte(specdata), nil))
		Expect(reflect.TypeOf(repo).String()).To(Equal("*npm.Repository"))
	})

	It("retrieves credentials", func() {
		repo := Must(DefaultContext.RepositoryForConfig([]byte(specdata), nil))

		creds := Must(repo.LookupCredentials("registry.npmjs.org"))
		Expect(creds.Properties()).To(Equal(props))

		creds = Must(repo.LookupCredentials("npm.registry.acme.com/api/npm"))
		Expect(creds.Properties()).To(Equal(props2))
	})

	It("can access the default context", func() {
		ctx := credentials.New()

		r := runtimefinalizer.GetRuntimeFinalizationRecorder(ctx)
		Expect(r).NotTo(BeNil())

		Must(ctx.RepositoryForConfig([]byte(specdata), nil))

		ci := cpi.NewConsumerIdentity(identity.CONSUMER_TYPE)
		Expect(ci).NotTo(BeNil())
		credentials := Must(cpi.CredentialsForConsumer(ctx.CredentialsContext(), ci))
		Expect(credentials).NotTo(BeNil())
		Expect(credentials.Properties()).To(Equal(props))
	})
})
