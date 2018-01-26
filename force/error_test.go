package force

import (
	"fmt"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("was not found", func() {

		apiErr := ApiErrors{
			&ApiError{ErrorCode: "NOT_FOUND"},
		}

		found, err := WasNotFound(apiErr)
		if err != nil {
			fmt.Println(err)
			GinkgoT().Error("expected WasNotFound not to return an error")
		}
		if !found {
			GinkgoT().Error("expected err to say it was not found")
		}

		apiErr = ApiErrors{
			&ApiError{ErrorCode: "SO_TOTALLY_FOUND"},
		}

		found, err = WasNotFound(apiErr)
		if err != nil {
			fmt.Println(err)
			GinkgoT().Error("expected WasNotFound not to return an error")
		}
		if found {
			GinkgoT().Error("expected err not to say it was not found")
		}
	})
})
