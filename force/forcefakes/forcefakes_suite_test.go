package forcefakes_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestForcefakes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Forcefakes Suite")
}
