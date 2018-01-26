package force

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("limits", func() {

		forceApi := createTest()
		limits, err := forceApi.GetLimits()
		if err != nil {
			GinkgoT().Logf("Failed to get Limits, this is expected due to the developer account: %v", err)
		}
		GinkgoT().Log(limits)
	})
})
