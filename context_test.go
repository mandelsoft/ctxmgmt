package ctxmgmt_test

import (
	me "github.com/mandelsoft/ctxmgmt/attributes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("area test", func() {
	It("can be garbage collected", func() {
		// ctxlog.Context().AddRule(logging.NewConditionRule(logging.DebugLevel, me.Realm))

		ctx := me.New()
		Expect(ctx.IsIdenticalTo(ctx)).To(BeTrue())

		ctx2 := ctx.AttributesContext()
		Expect(ctx.IsIdenticalTo(ctx2)).To(BeTrue())

		ctx3 := me.New()
		Expect(ctx.IsIdenticalTo(ctx3)).To(BeFalse())
	})
})
