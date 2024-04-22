package quota_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSpacequota(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Spacequota Suite")
}
