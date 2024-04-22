package ldap_integration_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLdap(t *testing.T) {
	RegisterFailHandler(Fail)
	if os.Getenv("RUN_LDAP_TESTS") == "" {
		t.Skip("skipping LDAP tests")
	}
	RunSpecs(t, "LDAP Suite")
}
