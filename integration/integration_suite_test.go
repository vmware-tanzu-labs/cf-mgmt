package integration_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	if os.Getenv("RUN_INTEGRATION_TESTS") == "" {
		t.Skip("skipping integration tests as need pcfdev running")
	}
	RunSpecs(t, "Integration Suite")
}
