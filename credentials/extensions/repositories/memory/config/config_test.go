package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/datacontext/config"
	"github.com/mandelsoft/datacontext/credentials"
	"github.com/mandelsoft/datacontext/credentials/extensions/repositories/memory"
	"github.com/mandelsoft/datacontext/utils"
)

var _ = Describe("configure credentials", func() {
	var ctx credentials.Context
	var cfg config.Context

	BeforeEach(func() {
		cfg = config.New()
		ctx = credentials.WithConfigs(cfg).New()
	})

	It("reads config with ref", func() {
		data, err := os.ReadFile("testdata/creds.yaml")
		Expect(err).To(Succeed())
		_, err = cfg.ApplyData(data, nil, "creds.yaml")
		Expect(err).To(Succeed())

		spec := memory.NewRepositorySpec("default")
		repo, err := ctx.RepositoryForSpec(spec)
		Expect(err).To(Succeed())
		mem := repo.(*memory.Repository)
		Expect(mem.ExistsCredentials("ref")).To(BeTrue())
		creds, err := mem.LookupCredentials("ref")
		Expect(err).To(Succeed())
		Expect(creds.Properties()).To(Equal(utils.Properties{"username": "mandelsoft", "password": "specialsecret"}))
	})

	It("reads config with direct", func() {
		data, err := os.ReadFile("testdata/creds.yaml")
		Expect(err).To(Succeed())
		_, err = cfg.ApplyData(data, nil, "creds.yaml")
		Expect(err).To(Succeed())

		spec := memory.NewRepositorySpec("default")
		repo, err := ctx.RepositoryForSpec(spec)
		Expect(err).To(Succeed())
		mem := repo.(*memory.Repository)
		Expect(mem.ExistsCredentials("direct")).To(BeTrue())
		creds, err := mem.LookupCredentials("direct")
		Expect(err).To(Succeed())
		Expect(creds.Properties()).To(Equal(utils.Properties{"username": "mandelsoft2", "password": "specialsecret2"}))
	})
})
