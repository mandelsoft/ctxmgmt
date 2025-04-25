package internal_test

import (
	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/attributes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	local "github.com/mandelsoft/ctxmgmt/config/internal"
)

var _ = Describe("builder test", func() {
	It("creates local", func() {
		ctx := local.Builder{}.New(ctxmgmt.MODE_SHARED)

		Expect(ctx.AttributesContext()).To(BeIdenticalTo(attributes.DefaultContext))
		Expect(ctx).NotTo(BeIdenticalTo(local.DefaultContext))
		Expect(ctx.ConfigTypes()).To(BeIdenticalTo(local.DefaultConfigTypeScheme))
	})

	It("creates configured", func() {
		ctx := local.Builder{}.New(ctxmgmt.MODE_CONFIGURED)

		Expect(ctx.AttributesContext()).NotTo(BeIdenticalTo(attributes.DefaultContext))
		Expect(ctx).NotTo(BeIdenticalTo(local.DefaultContext))
		Expect(ctx.ConfigTypes()).NotTo(BeIdenticalTo(local.DefaultConfigTypeScheme))
		Expect(ctx.ConfigTypes().KnownTypeNames()).To(Equal(local.DefaultConfigTypeScheme.KnownTypeNames()))
	})

	It("creates iniial", func() {
		ctx := local.Builder{}.New(ctxmgmt.MODE_INITIAL)

		Expect(ctx.AttributesContext()).NotTo(BeIdenticalTo(attributes.DefaultContext))
		Expect(ctx).NotTo(BeIdenticalTo(local.DefaultContext))
		Expect(ctx.ConfigTypes()).NotTo(BeIdenticalTo(local.DefaultConfigTypeScheme))
		Expect(len(ctx.ConfigTypes().KnownTypeNames())).To(Equal(0))
	})
})
