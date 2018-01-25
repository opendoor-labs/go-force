package force

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("oauth", func() {

		forceApi := createTest()

		if err := forceApi.oauth.Validate(); err != nil {
			GinkgoT().Fatalf("Oauth object is invlaid: %#v", err)
		}
	})
})
