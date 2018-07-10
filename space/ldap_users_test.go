package space_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	l "github.com/pivotalservices/cf-mgmt/ldap"
	ldap "github.com/pivotalservices/cf-mgmt/ldap/fakes"

	. "github.com/pivotalservices/cf-mgmt/space"
)

var _ = Describe("given UserManager", func() {
	var (
		mockLdap    *ldap.FakeManager
		userManager UserManager
	)

	BeforeEach(func() {
		mockLdap = new(ldap.FakeManager)
		userManager = UserManager{
			LdapMgr: mockLdap,
		}
	})

	Context("GetLdapUsers()", func() {
		It("update ldap group users where users are not in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}

			mockLdap.GetUserIDsReturns([]l.User{
				l.User{
					UserDN: "CN=Cwashburn,OU=User,OU=Users,DC=company,DC=com",
					UserID: "cwashburn@company.com",
					Email:  "cwashburn@company.com",
				},
				l.User{
					UserDN: "CN=Cwashburn2,OU=User,OU=Users,DC=company,DC=com",
					UserID: "cwashburn2@company.com",
					Email:  "cwashburn2@company.com",
				},
				l.User{
					UserDN: "CN=Cwashburn3,OU=User,OU=Users,DC=company,DC=com",
					UserID: "cwashburn3@company.com",
					Email:  "cwashburn3@company.com",
				},
			}, nil)
			userList, err := userManager.GetLdapUsers(config, []string{"test1", "test2"}, []string{"userList1", "userList2"})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(userList).ShouldNot(BeNil())
			Expect(len(userList)).Should(Equal(3))
		})
	})
})
