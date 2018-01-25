package sobjects_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSobjects(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sobjects Suite")
}
