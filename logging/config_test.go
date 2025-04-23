package logging_test

import (
	"bytes"

	. "github.com/mandelsoft/datacontext/logging/testhelper"
	"github.com/mandelsoft/goutils/substutils"
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/logging"
	logcfg "github.com/mandelsoft/logging/config"
	"github.com/tonglil/buflogr"

	local "github.com/mandelsoft/datacontext/logging"
)

////////////////////////////////////////////////////////////////////////////////

var _ = Describe("logging configuration", func() {
	var buf bytes.Buffer
	var ctx logging.Context

	BeforeEach(func() {
		local.SetContext(logging.NewDefault())
		buf.Reset()
		def := buflogr.NewWithBuffer(&buf)
		ctx = local.Context()
		ctx.SetBaseLogger(def)
	})

	It("just logs with defaults", func() {
		LogTest(ctx)
		Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm contexts
`, substutils.SubstList("REALM", local.REALM.Name())))
	})
	It("just logs with config", func() {
		r := logcfg.ConditionalRule("debug")
		cfg := &logcfg.Config{
			Rules: []logcfg.Rule{r},
		}

		Expect(local.Configure(cfg)).To(Succeed())
		LogTest(ctx)
		Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[4] debug realm ${REALM}
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm ${REALM}
`, substutils.SubstList("REALM", local.REALM.Name())))
	})
})
