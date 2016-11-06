package neutron_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNeutron(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Neutron Suite")
}
