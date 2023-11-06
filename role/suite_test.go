package role_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var test *testing.T

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	test = t
	RunSpecs(t, "Test Suite")
}
