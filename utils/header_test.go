package utils

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	dummyApiKey = "QUBIC_API_KEY"
)

var _ = Describe("Headers", func() {
	It("signature", func() {
		Expect(signature(
			"secret",
			"1566549227549",
			http.MethodPut,
			"/test/path?currency=USD",
			"the_body",
		)).Should(Equal("xN/7FHzMvIVbJYESYPJlMwNHL9r3DBZ21lsjSn5W3Bo="))
	})
})

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Header Signature Suite")
}
