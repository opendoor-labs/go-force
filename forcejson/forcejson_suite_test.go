package forcejson_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestForcejson(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Forcejson Suite")
}
