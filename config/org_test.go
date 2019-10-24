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

	Context("Should merge protected org list with default protected orgs", func() {
		It("should not include duplicates", func() {
			org := &Orgs{
				ProtectedOrgs: []string{"p-dataflow", "protect-me"},
			}
			protectedOrgList := org.ProtectedOrgList()
			Expect(protectedOrgList).Should(HaveLen(8))
			Expect(protectedOrgList).Should(ContainElement("system"))
			Expect(protectedOrgList).Should(ContainElement("p-spring-cloud-services"))
			Expect(protectedOrgList).Should(ContainElement("splunk-nozzle-org"))
			Expect(protectedOrgList).Should(ContainElement("redis-test-ORG*"))
			Expect(protectedOrgList).Should(ContainElement("appdynamics-org"))
			Expect(protectedOrgList).Should(ContainElement("credhub-service-broker-org"))
			Expect(protectedOrgList).Should(ContainElement("p-dataflow"))
			Expect(protectedOrgList).Should(ContainElement("protect-me"))
		})
	})
})
