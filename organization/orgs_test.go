package organization_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cf-mgmt/organization"
)

var _ = Describe("given OrgManager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewManager("test.com", "token", "uaacToken")
			Ω(manager).ShouldNot(BeNil())
			orgManager := manager.(*DefaultOrgManager)
			Ω(orgManager.Host).Should(Equal("https://api.test.com"))
		})
	})
})
