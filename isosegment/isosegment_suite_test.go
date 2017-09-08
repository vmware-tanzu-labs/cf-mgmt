package isosegment

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var testingT *testing.T

func TestIsosegment(t *testing.T) {
	testingT = t
	RegisterFailHandler(Fail)
	RunSpecs(t, "Isosegment Suite")
}
