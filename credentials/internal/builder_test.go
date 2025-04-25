package internal_test

import (
	"github.com/mandelsoft/ctxmgmt/attributes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/config"
	local "github.com/mandelsoft/ctxmgmt/credentials/internal"
)

var _ = Describe("builder test", func() {
	It("creates local", func() {
		ctx := local.Builder{}.New(ctxmgmt.MODE_SHARED)

		Expect(ctx.AttributesContext()).To(BeIdenticalTo(attributes.DefaultContext))
		Expect(ctx).NotTo(BeIdenticalTo(local.DefaultContext))
		Expect(ctx.RepositoryTypes()).To(BeIdenticalTo(local.DefaultRepositoryTypeScheme))

		Expect(ctx.ConfigContext().GetId()).To(BeIdenticalTo(config.DefaultContext().GetId()))
	})

	It("creates defaulted", func() {
		ctx := local.Builder{}.New(ctxmgmt.MODE_DEFAULTED)

		Expect(ctx.AttributesContext()).NotTo(BeIdenticalTo(attributes.DefaultContext))
		Expect(ctx).NotTo(BeIdenticalTo(local.DefaultContext))
		Expect(ctx.RepositoryTypes()).To(BeIdenticalTo(local.DefaultRepositoryTypeScheme))

		Expect(ctx.ConfigContext().GetId()).NotTo(BeIdenticalTo(config.DefaultContext().GetType()))
		Expect(ctx.ConfigContext().ConfigTypes()).To(BeIdenticalTo(config.DefaultContext().ConfigTypes()))
	})

	It("creates configured", func() {
		ctx := local.Builder{}.New(ctxmgmt.MODE_CONFIGURED)

		Expect(ctx.AttributesContext()).NotTo(BeIdenticalTo(attributes.DefaultContext))
		Expect(ctx).NotTo(BeIdenticalTo(local.DefaultContext))
		Expect(ctx.RepositoryTypes()).NotTo(BeIdenticalTo(local.DefaultRepositoryTypeScheme))
		Expect(ctx.RepositoryTypes().KnownTypeNames()).To(Equal(local.DefaultRepositoryTypeScheme.KnownTypeNames()))

		Expect(ctx.ConfigContext().GetId()).NotTo(BeIdenticalTo(config.DefaultContext().GetId()))
		Expect(ctx.ConfigContext().ConfigTypes()).NotTo(BeIdenticalTo(config.DefaultContext().ConfigTypes()))
		Expect(ctx.ConfigContext().ConfigTypes().KnownTypeNames()).To(Equal(config.DefaultContext().ConfigTypes().KnownTypeNames()))
	})

	It("creates iniial", func() {
		ctx := local.Builder{}.New(ctxmgmt.MODE_INITIAL)

		Expect(ctx.AttributesContext()).NotTo(BeIdenticalTo(attributes.DefaultContext))
		Expect(ctx).NotTo(BeIdenticalTo(local.DefaultContext))
		Expect(ctx.RepositoryTypes()).NotTo(BeIdenticalTo(local.DefaultRepositoryTypeScheme))
		Expect(len(ctx.RepositoryTypes().KnownTypeNames())).To(Equal(0))

		Expect(ctx.ConfigContext()).NotTo(BeIdenticalTo(config.DefaultContext()))
		Expect(len(ctx.ConfigContext().ConfigTypes().KnownTypeNames())).To(Equal(0))
	})
})
