package organizationreader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOrganizationreader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Organizationreader Suite")
}
