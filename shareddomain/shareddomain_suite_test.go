package shareddomain_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestShareddomain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shareddomain Suite")
}
