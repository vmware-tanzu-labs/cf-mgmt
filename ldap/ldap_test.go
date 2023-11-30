package ldap_test

import (
	"errors"

	l "github.com/go-ldap/ldap/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap/fakes"
)

var _ = Describe("Ldap", func() {
	Describe("given a ldap manager", func() {
		var ldapManager *ldap.Manager
		var connection *fakes.FakeConnection
		var ldapConfig *config.LdapConfig
		BeforeEach(func() {
			ldapConfig = &config.LdapConfig{
				Origin:            "ldap",
				GroupAttribute:    "member",
				UserNameAttribute: "uid",
				UserMailAttribute: "mail",
				GroupSearchBase:   "ou=groups,dc=pivotal,dc=org",
				UserSearchBase:    "ou=users,dc=pivotal,dc=org",
			}
			connection = &fakes.FakeConnection{}
			ldapManager = &ldap.Manager{Config: ldapConfig, Connection: connection}
		})
		Context("GetUserByID()", func() {
			It("should return specified user", func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{
						{
							DN: "cn=cwashburn,ou=users,dc=pivotal,dc=org",
							Attributes: []*l.EntryAttribute{
								{Name: "mail", Values: []string{"cwashburn@foo.com"}},
							}},
					},
				}, nil)
				user, err := ldapManager.GetUserByID("cwashburn")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(user).ShouldNot(BeNil())
				Expect(user.Email).Should(BeEquivalentTo("cwashburn@foo.com"))
				Expect(user.UserID).Should(BeEquivalentTo("cwashburn"))
				Expect(user.UserDN).Should(BeEquivalentTo("cn=cwashburn,ou=users,dc=pivotal,dc=org"))
			})

			It("should return nil user when multiple entries found", func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{
						{
							DN: "cn=cwashburn,ou=users,dc=pivotal,dc=org",
							Attributes: []*l.EntryAttribute{
								{Name: "mail", Values: []string{"cwashburn@foo.com"}},
							}},
						{
							DN: "cn=cwashburn,ou=users,dc=pivotal,dc=org",
							Attributes: []*l.EntryAttribute{
								{Name: "mail", Values: []string{"cwashburn@foo.com"}},
							}},
					},
				}, nil)
				user, err := ldapManager.GetUserByID("cwashburn")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(user).Should(BeNil())
			})

			It("should return error when search fails", func() {
				connection.SearchReturns(nil, errors.New("Error searching"))
				_, err := ldapManager.GetUserByID("cwashburn")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(BeEquivalentTo("Error searching"))
			})
		})

		Context("GetUserByDN()", func() {
			It("should return specified user", func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{
						{
							DN: "cn=cwashburn,ou=users,dc=pivotal,dc=org",
							Attributes: []*l.EntryAttribute{
								{Name: "mail", Values: []string{"cwashburn@foo.com"}},
								{Name: "uid", Values: []string{"cwashburn"}},
							}},
					},
				}, nil)
				user, err := ldapManager.GetUserByDN("cn=cwashburn,ou=users,dc=pivotal,dc=org")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(user).ShouldNot(BeNil())
				Expect(user.Email).Should(BeEquivalentTo("cwashburn@foo.com"))
				Expect(user.UserID).Should(BeEquivalentTo("cwashburn"))
				Expect(user.UserDN).Should(BeEquivalentTo("cn=cwashburn,ou=users,dc=pivotal,dc=org"))

				searchRequest := connection.SearchArgsForCall(0)
				Expect(searchRequest.Filter).Should(BeEquivalentTo("(cn=cwashburn)"))
			})

			It(`should return specified user when \, is in dn`, func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{
						{
							DN: `cn=Washburn\, Caleb,ou=users,dc=pivotal,dc=org`,
							Attributes: []*l.EntryAttribute{
								{Name: "mail", Values: []string{"cwashburn@foo.com"}},
								{Name: "uid", Values: []string{"cwashburn"}},
							}},
					},
				}, nil)
				user, err := ldapManager.GetUserByDN(`cn=Washburn\, Caleb,ou=users,dc=pivotal,dc=org`)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(user).ShouldNot(BeNil())
				Expect(user.Email).Should(BeEquivalentTo("cwashburn@foo.com"))
				Expect(user.UserID).Should(BeEquivalentTo("cwashburn"))
				Expect(user.UserDN).Should(BeEquivalentTo(`cn=Washburn\, Caleb,ou=users,dc=pivotal,dc=org`))

				searchRequest := connection.SearchArgsForCall(0)
				Expect(searchRequest.Filter).Should(BeEquivalentTo("(cn=Washburn, Caleb)"))
			})

			It("should return specified user when space is in dn", func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{
						{
							DN: "cn=Caleb A. Washburn,ou=users,dc=pivotal,dc=org",
							Attributes: []*l.EntryAttribute{
								{Name: "mail", Values: []string{"cwashburn@foo.com"}},
								{Name: "uid", Values: []string{"cwashburn"}},
							}},
					},
				}, nil)
				user, err := ldapManager.GetUserByDN("cn=Caleb A. Washburn,ou=users,dc=pivotal,dc=org")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(user).ShouldNot(BeNil())
				Expect(user.Email).Should(BeEquivalentTo("cwashburn@foo.com"))
				Expect(user.UserID).Should(BeEquivalentTo("cwashburn"))
				Expect(user.UserDN).Should(BeEquivalentTo("cn=Caleb A. Washburn,ou=users,dc=pivotal,dc=org"))

				searchRequest := connection.SearchArgsForCall(0)
				Expect(searchRequest.Filter).Should(BeEquivalentTo("(cn=Caleb A. Washburn)"))
			})

			It("should return nil user when multiple entries found", func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{
						{
							DN: "cn=cwashburn,ou=users,dc=pivotal,dc=org",
							Attributes: []*l.EntryAttribute{
								{Name: "mail", Values: []string{"cwashburn@foo.com"}},
								{Name: "uid", Values: []string{"cwashburn"}},
							}},
						{
							DN: "cn=cwashburn,ou=users,dc=pivotal,dc=org",
							Attributes: []*l.EntryAttribute{
								{Name: "mail", Values: []string{"cwashburn@foo.com"}},
								{Name: "uid", Values: []string{"cwashburn"}},
							}},
					},
				}, nil)
				user, err := ldapManager.GetUserByDN("cn=cwashburn,ou=users,dc=pivotal,dc=org")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(user).Should(BeNil())
			})

			It("should return error when search fails", func() {
				connection.SearchReturns(nil, errors.New("Error searching"))
				_, err := ldapManager.GetUserByDN("cn=cwashburn,ou=users,dc=pivotal,dc=org")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(BeEquivalentTo("Error searching"))
			})

			It("should return error when invalid cn", func() {
				_, err := ldapManager.GetUserByDN("cwashburn")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(BeEquivalentTo("cannot find CN for DN: cwashburn"))
				Expect(connection.SearchCallCount()).Should(Equal(0))
			})
		})
		Context("GetUserDNs()", func() {
			It("should return users for specified group", func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{
						{
							Attributes: []*l.EntryAttribute{
								{Name: "member", Values: []string{
									"cn=cwashburn,ou=users,dc=pivotal,dc=org",
									"cn=cwashburn1,ou=users,dc=pivotal,dc=org",
									`cn=Washburn\, Caleb,ou=users,dc=pivotal,dc=org`}},
							}},
					},
				}, nil)
				users, err := ldapManager.GetUserDNs("group1")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(users)).Should(Equal(3))
				Expect(users).Should(ConsistOf([]string{"cn=cwashburn,ou=users,dc=pivotal,dc=org", "cn=cwashburn1,ou=users,dc=pivotal,dc=org", `cn=Washburn\, Caleb,ou=users,dc=pivotal,dc=org`}))
			})

			It("should return empty list when group is not found", func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{},
				}, nil)
				users, err := ldapManager.GetUserDNs("group1")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(users)).Should(Equal(0))
			})

			It("should return empty list when group has no users", func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{
						{},
					},
				}, nil)
				users, err := ldapManager.GetUserDNs("group1")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(users)).Should(Equal(0))
			})
			It("should return empty list when multiple groups are found", func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{
						{},
						{},
					},
				}, nil)
				users, err := ldapManager.GetUserDNs("group1")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(users)).Should(Equal(0))
			})

			It("should return error when search fails", func() {
				connection.SearchReturns(nil, errors.New("Error searching"))
				_, err := ldapManager.GetUserDNs("group1")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(BeEquivalentTo("Error searching"))
			})
		})

		Context("GroupFilter()", func() {
			It("Should return expected group filter", func() {
				filter, err := ldapManager.GroupFilter("cn=nested_group,ou=groups,dc=pivotal,dc=org")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(filter).Should(Equal("(&(objectclass=groupOfNames)(cn=nested_group))"))
			})

			It("Should return expected group filter", func() {
				filter, err := ldapManager.GroupFilter("CN=nested_group,ou=groups,dc=pivotal,dc=org")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(filter).Should(Equal("(&(objectclass=groupOfNames)(CN=nested_group))"))
			})

			It("Should error", func() {
				_, err := ldapManager.GroupFilter("foo")
				Expect(err).Should(MatchError("cannot find CN for DN: foo"))
			})
		})

		Context("IsGroup()", func() {
			It("Should return false", func() {
				isGroup, groupName, err := ldapManager.IsGroup("foo")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(isGroup).Should(BeFalse())
				Expect(groupName).Should(Equal(""))
			})

			It("Should return true", func() {
				connection.SearchReturns(&l.SearchResult{
					Entries: []*l.Entry{
						{
							Attributes: []*l.EntryAttribute{
								{
									Name:   "cn",
									Values: []string{"nested_group"},
								},
							}},
					},
				}, nil)
				isGroup, groupName, err := ldapManager.IsGroup("cn=nested_group,ou=groups,dc=pivotal,dc=org")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(isGroup).Should(BeTrue())
				Expect(groupName).Should(Equal("nested_group"))
			})
		})
	})
})
