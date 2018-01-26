package force_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestForce(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Force Suite")
}
