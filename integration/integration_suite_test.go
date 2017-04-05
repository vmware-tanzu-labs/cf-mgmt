package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	if testing.Short() {
		t.Skip("skipping integration tests as need pcfdev running")
	}
	RunSpecs(t, "Integration Suite")
}
