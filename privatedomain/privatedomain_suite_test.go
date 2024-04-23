package privatedomain_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPrivatedomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Privatedomain Suite")
}
