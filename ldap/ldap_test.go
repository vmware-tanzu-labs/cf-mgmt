package ldap_test

import (
	. "github.com/pivotalservices/cf-mgmt/ldap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = XDescribe("Ldap", func() {
	var ldapManager Manager
	var mgrError error
	Describe("given a GetUserIDs", func() {
		BeforeEach(func() {
			if ldapManager, mgrError = NewDefaultManager("./fixtures/example"); mgrError != nil {
				panic(mgrError)
			}
		})
		Context("when called with a valid group", func() {
			It("then it should return 3 users", func() {
				users, err := ldapManager.GetUserIDs("test_space1_developers")
				Ω(err).Should(BeNil())
				Ω(len(users)).Should(Equal(3))
			})

		})
	})
})
