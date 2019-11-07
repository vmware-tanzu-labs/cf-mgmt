package organization_test

import (
	. "github.com/pivotalservices/cf-mgmt/organization"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrgMatcher", func() {
	Context("Matches", func() {
		It("Is false", func() {
			Expect(Matches("foo", []string{})).To(BeFalse())
		})
		It("Is True", func() {
			Expect(Matches("foo", []string{"foo"})).To(BeTrue())
		})
		It("Is True", func() {
			Expect(Matches("sandbox-org", []string{"org-*"})).To(BeTrue())
		})
		It("Is True", func() {
			Expect(Matches("redis-test-ORG1233", []string{"redis-test-ORG*"})).To(BeTrue())
		})
	})
})
