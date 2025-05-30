package config_test

import (
	"encoding/json"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/goutils/testutils"

	"github.com/mandelsoft/ctxmgmt/config"
	"github.com/mandelsoft/ctxmgmt/credentials"
	localconfig "github.com/mandelsoft/ctxmgmt/credentials/config"
	"github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/aliases"
	"github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/directcreds"
	"github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/memory"
	"github.com/mandelsoft/ctxmgmt/utils"
	"github.com/mandelsoft/ctxmgmt/utils/runtime"
)

var DefaultContext = credentials.New()

var _ = Describe("generic credentials", func() {
	props := utils.Properties{
		"user":     "USER",
		"password": "PASSWORD",
	}

	repospec := memory.NewRepositorySpec("test")
	credspec := credentials.NewCredentialsSpec("cred", repospec)
	direct := directcreds.NewRepositorySpec(props)

	cfgconsumerdata := "{\"type\":\"credentials.config.mandelsoft.de\",\"consumers\":[{\"identity\":{\"type\":\"oci\",\"url\":\"https://acme.com\"},\"credentials\":[{\"credentialsName\":\"cred\",\"repoName\":\"test\",\"type\":\"Memory\"}]}]}"
	cfgrepodata := "{\"type\":\"credentials.config.mandelsoft.de\",\"repositories\":[{\"repository\":{\"repoName\":\"test\",\"type\":\"Memory\"},\"credentials\":[{\"properties\":{\"password\":\"PASSWORD\",\"user\":\"USER\"},\"type\":\"Credentials\"}]}]}"
	cfgaliasdata := "{\"type\":\"credentials.config.mandelsoft.de\",\"aliases\":{\"alias\":{\"repository\":{\"repoName\":\"test\",\"type\":\"Memory\"},\"credentials\":[{\"properties\":{\"password\":\"PASSWORD\",\"user\":\"USER\"},\"type\":\"Credentials\"}]}}}"
	_ = props

	Context("serialize", func() {
		It("serializes repository spec not in map", func() {
			mapdata := "{\"repositories\":{\"repository\":{\"repoName\":\"test\",\"type\":\"Memory\"}}}"
			type S struct {
				Repositories localconfig.RepositorySpec `json:"repositories"`
			}

			rspec, err := credentials.ToGenericRepositorySpec(repospec)
			Expect(err).To(Succeed())
			s := &S{
				Repositories: localconfig.RepositorySpec{Repository: *rspec},
			}
			data, err := json.Marshal(s)

			Expect(err).To(Succeed())
			Expect(data).To(Equal([]byte(mapdata)))
		})

		It("serializes repository spec map", func() {
			mapdata := "{\"repositories\":{\"repo\":{\"repository\":{\"repoName\":\"test\",\"type\":\"Memory\"}}}}"
			type S struct {
				Repositories map[string]localconfig.RepositorySpec `json:"repositories"`
			}

			rspec, err := credentials.ToGenericRepositorySpec(repospec)
			Expect(err).To(Succeed())
			s := &S{
				Repositories: map[string]localconfig.RepositorySpec{
					"repo": {Repository: *rspec},
				},
			}
			data, err := json.Marshal(s)
			Expect(err).To(Succeed())
			Expect(data).To(Equal([]byte(mapdata)))
		})
	})

	Context("composition", func() {
		It("composes a config for consumers", func() {
			consumerid := credentials.ConsumerIdentity{
				"type": "oci",
				"url":  "https://acme.com",
			}

			cfg := localconfig.New()

			cfg.AddConsumer(consumerid, credspec)

			data, err := json.Marshal(cfg)
			Expect(err).To(Succeed())
			Expect(data).To(Equal([]byte(cfgconsumerdata)))

			cfg2 := &localconfig.Config{}
			err = json.Unmarshal(data, cfg2)
			Expect(err).To(Succeed())
			Expect(cfg2).To(Equal(cfg))
		})

		It("composes a config for repositories", func() {
			cfg := localconfig.New()

			cfg.AddRepository(repospec, direct)

			data, err := json.Marshal(cfg)
			Expect(err).To(Succeed())
			Expect(data).To(Equal([]byte(cfgrepodata)))

			cfg2 := &localconfig.Config{}
			err = json.Unmarshal(data, cfg2)
			Expect(err).To(Succeed())
			Expect(cfg2).To(Equal(cfg))
		})

		It("composes a config for aliases", func() {
			cfg := localconfig.New()

			cfg.AddAlias("alias", repospec, direct)

			data, err := json.Marshal(cfg)
			Expect(err).To(Succeed())
			Expect(data).To(Equal([]byte(cfgaliasdata)))

			cfg2 := &localconfig.Config{}
			err = json.Unmarshal(data, cfg2)
			Expect(err).To(Succeed())
			Expect(cfg2).To(Equal(cfg))
		})
	})

	Context("apply", func() {
		var ctx credentials.Context

		_ = ctx

		BeforeEach(func() {
			ctx = credentials.WithConfigs(config.New()).New()
		})

		It("applies a config for aliases", func() {
			cfg := localconfig.New()
			cfg.AddAlias("alias", repospec, direct)

			ctx.ConfigContext().ApplyConfig(cfg, "testconfig")

			spec := aliases.NewRepositorySpec("alias")

			repo, err := ctx.RepositoryForSpec(spec)
			Expect(err).To(Succeed())
			Expect(reflect.TypeOf(repo).String()).To(Equal("*memory.Repository"))
		})

		It("applies a config for consumers", func() {
			cfg := localconfig.New()

			consumer := credentials.ConsumerIdentity{
				credentials.ID_TYPE: "mytype",
				"host":              "localhost",
			}
			props := utils.Properties{"token": "mytoken"}
			creds := directcreds.NewCredentials(props)
			Expect(cfg.AddConsumer(consumer, creds)).To(Succeed())

			data, err := runtime.DefaultYAMLEncoding.Marshal(cfg)
			Expect(err).To(Succeed())
			Expect(string(data)).To(testutils.StringEqualTrimmedWithContext(`
consumers:
- credentials:
  - credentialsName: Credentials
    properties:
      token: mytoken
    type: Credentials
  identity:
    host: localhost
    type: mytype
type: credentials.config.mandelsoft.de
`))

			ctx.ConfigContext().ApplyConfig(cfg, "testconfig")

			result, err := credentials.CredentialsForConsumer(ctx, consumer, credentials.CompleteMatch)
			Expect(err).To(Succeed())

			Expect(result.Properties()).To(Equal(props))
		})

		It("applies a config for consumers", func() {
			props := utils.Properties{"token": "mytoken"}
			consumer := credentials.ConsumerIdentity{
				credentials.ID_TYPE: "mytype",
				"host":              "localhost",
			}
			data := `
type: credentials.config.mandelsoft.de
consumers:
- credentials:
  - type: Credentials
    properties:
      token: mytoken
  identity:
    host: localhost
    type: mytype
`
			ctx.ConfigContext().ApplyData([]byte(data), nil, "testconfig")

			result, err := credentials.CredentialsForConsumer(ctx, consumer, credentials.CompleteMatch)
			Expect(err).To(Succeed())

			Expect(result.Properties()).To(Equal(props))
		})
	})
})
