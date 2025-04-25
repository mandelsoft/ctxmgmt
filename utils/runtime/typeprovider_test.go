package runtime_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/mandelsoft/ctxmgmt/utils/runtime"
)

var _ = Describe("Type Provider Test Environment", func() {
	Context("K8S", func() {
		It("maps base apiGroup", func() {
			t, ok := runtime.MapK8SManifestInfoToType("v1", "Secret")
			Expect(ok).To(BeTrue())
			Expect(t).To(Equal("Secret/v1"))
		})

		It("maps apiGroup", func() {
			t, ok := runtime.MapK8SManifestInfoToType("apps/v1", "Deployment")
			Expect(ok).To(BeTrue())
			Expect(t).To(Equal("Deployment.apps/v1"))
		})

		It("remaps base apiGroup", func() {
			g, k := runtime.MapTypeToK8SManifestInfo("Secret/v1")
			Expect(g).To(Equal("v1"))
			Expect(k).To(Equal("Secret"))
		})

		It("remaps apiGroup", func() {
			g, k := runtime.MapTypeToK8SManifestInfo("Deployment.apps/v1")
			Expect(g).To(Equal("apps/v1"))
			Expect(k).To(Equal("Deployment"))
		})

		It("extracts type", func() {
			t, ok := runtime.DefaultTypeProviderRegistry.GetTypeFor([]byte(`
apiVersion: test/v1
kind: object
`), runtime.DefaultYAMLEncoding)
			Expect(ok).To(BeTrue())
			Expect(t).To(Equal("object.test/v1"))
		})
	})

	Context("default", func() {
		It("extracts type", func() {
			t, ok := runtime.DefaultTypeProviderRegistry.GetTypeFor([]byte(`
type: test/v1
`), runtime.DefaultYAMLEncoding)
			Expect(ok).To(BeTrue())
			Expect(t).To(Equal("test/v1"))
		})
	})
})
