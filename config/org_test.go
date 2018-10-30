package config_test

import (
	. "github.com/pivotalservices/cf-mgmt/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Org", func() {
	Context("Replacing Org in list", func() {
		It("Should not contain old name and should contain new name", func() {
			org := &Orgs{
				Orgs: []string{"foo", "bar"},
			}
			org.Replace("foo", "new-foo")
			Expect(org.Orgs).To(ConsistOf([]string{"new-foo", "bar"}))
			Expect(len(org.Orgs)).To(Equal(2))
		})
	})
})
