package api_test

import (
	"github.com/mandelsoft/datacontext/attributes"
	. "github.com/mandelsoft/goutils/testutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/datacontext"
	"github.com/mandelsoft/datacontext/action/api"
	"github.com/mandelsoft/datacontext/action/handlers"
	"github.com/mandelsoft/datacontext/utils"
	"github.com/mandelsoft/datacontext/utils/runtime"
)

type Handler struct {
	spec  api.ActionSpec
	creds utils.Properties
}

func (h *Handler) Handle(spec api.ActionSpec, creds utils.Properties) (api.ActionResult, error) {
	h.spec = spec
	h.creds = creds
	r := NewActionResult(spec.(*ActionSpec).Field)
	r.SetVersion("v2")
	return r, nil
}

var _ handlers.ActionHandler = &Handler{}

var _ = Describe("action registry", func() {
	var registry api.ActionTypeRegistry

	BeforeEach(func() {
		registry = api.NewActionTypeRegistry()
		RegisterAction(registry)
	})

	Context("plain", func() {
		It("registers", func() {
			Expect(registry.SupportedActionVersions(NAME)).To(Equal([]string{"v1", "v2"}))
		})

		It("encoding spec v1", func() {
			spec := NewActionSpec("acme.com")
			spec.SetVersion("v1")
			data := Must(registry.EncodeActionSpec(spec, runtime.DefaultJSONEncoding))
			Expect(string(data)).To(Equal(`{"type":"testAction/v1","field":"acme.com"}`))
			d := Must(registry.DecodeActionSpec(data, runtime.DefaultJSONEncoding))
			Expect(d).To(Equal(spec))
		})
		It("encoding spec v2", func() {
			spec := NewActionSpec("acme.com")
			spec.SetVersion("v2")
			data := Must(registry.EncodeActionSpec(spec, runtime.DefaultJSONEncoding))
			Expect(string(data)).To(Equal(`{"type":"testAction/v2","data":"acme.com"}`))
			d := Must(registry.DecodeActionSpec(data, runtime.DefaultJSONEncoding))
			Expect(d).To(Equal(spec))
		})

		It("encoding result v1", func() {
			spec := NewActionResult("successful")
			spec.SetVersion("v1")
			data := Must(registry.EncodeActionResult(spec, runtime.DefaultJSONEncoding))
			Expect(string(data)).To(Equal(`{"type":"testAction/v1","message":"successful"}`))
			d := Must(registry.DecodeActionResult(data, runtime.DefaultJSONEncoding))
			Expect(d).To(Equal(spec))
		})
		It("encoding result v2", func() {
			spec := NewActionResult("successful")
			spec.SetVersion("v2")
			data := Must(registry.EncodeActionResult(spec, runtime.DefaultJSONEncoding))
			Expect(string(data)).To(Equal(`{"type":"testAction/v2","data":"successful"}`))
			d := Must(registry.DecodeActionResult(data, runtime.DefaultJSONEncoding))
			Expect(d).To(Equal(spec))
		})
	})

	Context("data context", func() {
		var ctx datacontext.Context
		var handler *Handler

		BeforeEach(func() {
			handler = &Handler{}
			ctx = attributes.NewWithActions(nil, handlers.NewRegistry(registry))
			Expect(ctx.GetActions().GetActionTypes()).To(BeIdenticalTo(registry))
			MustBeSuccessful(ctx.GetActions().Register(handler, handlers.ForAction(NAME), handlers.WithVersions("v2"), handlers.ForSelectors(".*\\.com")))
		})

		It("", func() {
			spec := NewActionSpec("acme.com")
			creds := utils.Properties{"alice": "bob"}
			r := Must(ctx.GetActions().Execute(spec, creds))
			Expect(handler.spec).To(Equal(spec))
			Expect(handler.creds).To(Equal(creds))
			rs := NewActionResult("acme.com")
			rs.SetVersion("v2")
			Expect(r).To(Equal(rs))
		})
	})
})
