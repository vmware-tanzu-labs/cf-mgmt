package uaa_test

import (
	"errors"

	uaaclient "github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/vmwarepivotallabs/cf-mgmt/uaa"

	"github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
)

var _ = Describe("given uaa manager", func() {
	var (
		fakeuaa *fakes.FakeUaa
		manager DefaultUAAManager
	)
	BeforeEach(func() {
		fakeuaa = &fakes.FakeUaa{}
		manager = DefaultUAAManager{
			Client: fakeuaa,
		}
	})

	Context("ListUsers()", func() {

		It("should return list of users", func() {
			fakeuaa.ListAllUsersReturns([]uaaclient.User{
				{Username: "foo4", ID: "foo4-id"},
				{Username: "admin", ID: "admin-id"},
				{Username: "user", ID: "user-id"},
				{Username: "cwashburn", ID: "cwashburn-id"},
				{Username: "foo", ID: "foo-id"},
				{Username: "foo1", ID: "foo1-id"},
				{Username: "foo2", ID: "foo2-id"},
				{Username: "foo3", ID: "foo3-id"},
				{Username: "cn=admin", ID: "cn=admin-id"},
			}, nil)
			users, err := manager.ListUsers()
			Ω(err).ShouldNot(HaveOccurred())
			keys := make([]string, 0, len(users.List()))
			for _, k := range users.List() {
				keys = append(keys, k.Username)
			}
			Ω(len(users.List())).Should(Equal(9))
			Ω(keys).Should(ConsistOf("foo4", "admin", "user", "cwashburn", "foo", "foo1", "foo2", "foo3", "cn=admin"))
		})
		It("should return an error", func() {
			fakeuaa.ListAllUsersReturns(nil, errors.New("Got an error"))
			_, err := manager.ListUsers()
			Ω(err).Should(HaveOccurred())
		})
	})
	Context("CreateLdapUser()", func() {

		It("should successfully create user", func() {
			userName := "user"
			userEmail := "email"
			externalID := "userDN"

			fakeuaa.CreateUserReturns(
				&uaaclient.User{
					Username:   userName,
					ExternalID: externalID,
					Emails: []uaaclient.Email{
						{Value: userEmail},
					}},
				nil,
			)
			_, err := manager.CreateExternalUser(userName, userEmail, externalID, "ldap")
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should successfully create user with complex dn", func() {
			userName := "asdfasdfsadf"
			userEmail := "caleb.washburn@test.com"
			externalID := `CN=Washburn\, Caleb\, asdfasdfsadf\,OU=NO-HOME-USERS,OU=BU-USA,DC=1DC,DC=com`

			fakeuaa.CreateUserReturns(
				&uaaclient.User{
					Username:   userName,
					ExternalID: externalID,
					Emails: []uaaclient.Email{
						{Value: userEmail},
					}},
				nil,
			)
			_, err := manager.CreateExternalUser(userName, userEmail, externalID, "ldap")
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should peek", func() {
			userName := "user"
			userEmail := "email"
			externalID := "userDN"
			manager.Peek = true
			_, err := manager.CreateExternalUser(userName, userEmail, externalID, "ldap")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(fakeuaa.CreateUserCallCount()).Should(Equal(0))
		})
		It("should not invoke post", func() {
			_, err := manager.CreateExternalUser("", "", "", "ldap")
			Ω(err).Should(HaveOccurred())
			Ω(fakeuaa.CreateUserCallCount()).Should(Equal(0))
		})
	})
	Context("CreateSamlUser()", func() {
		It("should successfully create user", func() {
			userName := "user@test.com"
			userEmail := "user@test.com"
			externalID := "user@test.com"
			origin := "saml"

			fakeuaa.CreateUserReturns(
				&uaaclient.User{
					Username:   userName,
					ExternalID: externalID,
					Origin:     origin,
					Emails: []uaaclient.Email{
						{Value: userEmail},
					}},
				nil,
			)
			_, err := manager.CreateExternalUser(userName, userEmail, externalID, origin)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})
})
