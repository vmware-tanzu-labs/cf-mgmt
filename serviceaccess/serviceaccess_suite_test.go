package serviceaccess_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServiceaccess(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Serviceaccess Suite")
}
