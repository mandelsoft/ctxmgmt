package logging_test

import (
	"bytes"

	"github.com/mandelsoft/datacontext/attributes"
	. "github.com/mandelsoft/datacontext/logging/testhelper"
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/goutils/substutils"
	"github.com/mandelsoft/logging"
	"github.com/tonglil/buflogr"

	logcfg "github.com/mandelsoft/datacontext/attributes/config/logging"
	"github.com/mandelsoft/datacontext/config"
	ctxlog "github.com/mandelsoft/datacontext/logging"
)

var _ = Describe("logging configuration", func() {
	var ctx attributes.AttributesContext
	var cfg config.Context
	var buf bytes.Buffer
	var orig logging.Context

	BeforeEach(func() {
		orig = logging.DefaultContext().(*logging.ContextReference).Context
		logging.SetDefaultContext(logging.NewDefault())
		ctxlog.SetContext(nil)
		ctx = attributes.New(nil)
		cfg = config.WithSharedAttributes(ctx).New()

		buf.Reset()
		def := buflogr.NewWithBuffer(&buf)
		ctx.LoggingContext().SetBaseLogger(def)
	})

	AfterEach(func() {
		// logging.SetDefaultContext(orig)
	})
	_ = cfg
	_ = orig

	It("just logs with defaults", func() {
		LogTest(ctx)

		Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm ${REALM}
`, substutils.SubstList("REALM", ctxlog.REALM.Name())))
	})

	It("just logs with settings from default context", func() {
		logging.DefaultContext().AddRule(logging.NewConditionRule(logging.DebugLevel))
		LogTest(ctx)

		Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[4] debug realm ${REALM}
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm ${REALM}
`, substutils.SubstList("REALM", ctxlog.REALM.Name())))
	})

	It("just logs with settings from default context", func() {
		logging.DefaultContext().AddRule(logging.NewConditionRule(logging.DebugLevel))
		LogTest(cfg)

		Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[4] debug realm ${REALM}
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm ${REALM}
`, substutils.SubstList("REALM", ctxlog.REALM.Name())))
	})

	It("just logs with settings for root context", func() {
		spec := `
type: ` + logcfg.ConfigTypeV1 + `
contextType: ` + attributes.CONTEXT_TYPE + `
settings:
  rules:
  - rule:
      level: Debug
`
		_, err := cfg.ApplyData([]byte(spec), nil, "testconfig")
		Expect(err).To(Succeed())
		LogTest(ctx)
		LogTest(cfg, "cfg")

		Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[4] debug realm ${REALM}
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm ${REALM}
V[4] cfgdebug realm ${REALM}
V[3] cfginfo realm ${REALM}
V[2] cfgwarn realm ${REALM}
ERROR <nil> cfgerror realm ${REALM}
`, substutils.SubstList("REALM", ctxlog.REALM.Name())))
	})

	It("just logs with settings for root context by context provider", func() {
		spec := `
type: ` + logcfg.ConfigTypeV1 + `
settings:
  rules:
  - rule:
      level: Debug
`
		_, err := cfg.ApplyData([]byte(spec), nil, "testconfig")
		Expect(err).To(Succeed())

		LogTest(ctx)
		LogTest(cfg, "cfg")

		Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[4] debug realm ${REALM}
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm ${REALM}
V[4] cfgdebug realm ${REALM}
V[3] cfginfo realm ${REALM}
V[2] cfgwarn realm ${REALM}
ERROR <nil> cfgerror realm ${REALM}
`, substutils.SubstList("REALM", ctxlog.REALM.Name())))
	})

	It("just logs with settings for config context", func() {
		spec := `
type: ` + logcfg.ConfigTypeV1 + `
contextType: ` + config.CONTEXT_TYPE + `
settings:
  rules:
  - rule:
      level: Debug
`
		_, err := cfg.ApplyData([]byte(spec), nil, "testconfig")
		Expect(err).To(Succeed())

		LogTest(ctx)
		LogTest(cfg, "cfg")

		Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm ${REALM}
V[4] cfgdebug realm ${REALM}
V[3] cfginfo realm ${REALM}
V[2] cfgwarn realm ${REALM}
ERROR <nil> cfgerror realm ${REALM}
`, substutils.SubstList("REALM", ctxlog.REALM.Name())))
	})

	Context("default logging", func() {
		spec1 := `
type: ` + logcfg.ConfigTypeV1 + `
contextType: default
settings:
  rules:
  - rule:
      level: Debug
`
		spec2 := `
type: ` + logcfg.ConfigTypeV1 + `
contextType: default
settings:
  rules:
  - rule:
      level: Info
`
		spec3 := `
type: ` + logcfg.ConfigTypeV1 + `
contextType: default
extraId: extra
settings:
  rules:
  - rule:
      level: Debug
`

		var ctx logging.Context

		BeforeEach(func() {
			ctxlog.SetContext(logging.NewDefault())
			buf.Reset()
			def := buflogr.NewWithBuffer(&buf)
			ctx = ctxlog.Context()
			ctx.SetBaseLogger(def)
		})

		It("just logs with config", func() {
			_, err := cfg.ApplyData([]byte(spec1), nil, "spec1")
			Expect(err).To(Succeed())
			LogTest(ctx)

			Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[4] debug realm ${REALM}
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm ${REALM}
`, substutils.SubstList("REALM", ctxlog.REALM.Name())))
		})

		It("applies config once", func() {
			_, err := cfg.ApplyData([]byte(spec1), nil, "spec1")
			Expect(err).To(Succeed())
			_, err = cfg.ApplyData([]byte(spec2), nil, "spec2")
			Expect(err).To(Succeed())
			_, err = cfg.ApplyData([]byte(spec1), nil, "spec1.2")
			Expect(err).To(Succeed())

			LogTest(ctx)

			Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm ${REALM}
`, substutils.SubstList("REALM", ctxlog.REALM.Name())))
		})
		It("re-applies config with extra id", func() {
			_, err := cfg.ApplyData([]byte(spec1), nil, "spec1")
			Expect(err).To(Succeed())
			_, err = cfg.ApplyData([]byte(spec2), nil, "spec2")
			Expect(err).To(Succeed())
			_, err = cfg.ApplyData([]byte(spec3), nil, "spec3")
			Expect(err).To(Succeed())

			LogTest(ctx)

			Expect(buf.String()).To(StringEqualTrimmedWithContext(`
V[4] debug realm ${REALM}
V[3] info realm ${REALM}
V[2] warn realm ${REALM}
ERROR <nil> error realm ${REALM}
`, substutils.SubstList("REALM", ctxlog.REALM.Name())))
		})
	})
})
