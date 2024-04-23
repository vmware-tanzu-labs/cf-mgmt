package uaa_test

import (
	"errors"

	uaaclient "github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/vmwarepivotallabs/cf-mgmt/uaa"

	"github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
)

type UaaResponse struct {
	Response []uaaclient.User `json:"resources"`
}

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
			fakeuaa.ListUsersReturns([]uaaclient.User{
				{Username: "foo4", ID: "foo4-id"},
				{Username: "admin", ID: "admin-id"},
				{Username: "user", ID: "user-id"},
				{Username: "cwashburn", ID: "cwashburn-id"},
				{Username: "foo", ID: "foo-id"},
				{Username: "foo1", ID: "foo1-id"},
				{Username: "foo2", ID: "foo2-id"},
				{Username: "foo3", ID: "foo3-id"},
				{Username: "cn=admin", ID: "cn=admin-id"},
			}, uaaclient.Page{ItemsPerPage: 500, StartIndex: 1, TotalResults: 9}, nil)
			users, err := manager.ListUsers()
			Expect(fakeuaa.ListUsersCallCount()).Should(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			keys := make([]string, 0, len(users.List()))
			for _, k := range users.List() {
				keys = append(keys, k.Username)
			}
			Expect(len(users.List())).Should(Equal(9))
			Expect(keys).Should(ConsistOf("foo4", "admin", "user", "cwashburn", "foo", "foo1", "foo2", "foo3", "cn=admin"))
		})
		It("should return an error", func() {
			fakeuaa.ListUsersReturns(nil, uaaclient.Page{ItemsPerPage: 500, StartIndex: 1, TotalResults: 10}, errors.New("Got an error"))
			_, err := manager.ListUsers()
			Expect(err).Should(HaveOccurred())
			Expect(fakeuaa.ListUsersCallCount()).Should(Equal(1))
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
			err := manager.CreateExternalUser(userName, userEmail, externalID, "ldap")
			Expect(err).ShouldNot(HaveOccurred())
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
			err := manager.CreateExternalUser(userName, userEmail, externalID, "ldap")
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should peek", func() {
			userName := "user"
			userEmail := "email"
			externalID := "userDN"
			manager.Peek = true
			err := manager.CreateExternalUser(userName, userEmail, externalID, "ldap")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeuaa.CreateUserCallCount()).Should(Equal(0))
		})
		It("should not invoke post", func() {
			err := manager.CreateExternalUser("", "", "", "ldap")
			Expect(err).Should(HaveOccurred())
			Expect(fakeuaa.CreateUserCallCount()).Should(Equal(0))
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
			err := manager.CreateExternalUser(userName, userEmail, externalID, origin)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
