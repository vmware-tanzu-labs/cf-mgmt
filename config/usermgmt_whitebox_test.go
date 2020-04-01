package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserMgmt Whitebox Tests", func() {
	const (
		username = "my-user"
	)

	var u UserMgmt
	BeforeEach(func() {
		u = UserMgmt{}
	})

	Context("A user is requested to be added to the UserMgmt", func() {
		When("that user does not already exist", func() {
			When("the origin is internal", func() {
				It("adds the user", func() {
					u.addUser(InternalOrigin, username)
					Expect(u.Users).To(HaveLen(1))
					Expect(u.Users[0]).To(Equal(username))
					Expect(u.SamlUsers).To(HaveLen(0))
					Expect(u.LDAPUsers).To(HaveLen(0))
				})
			})

			When("the origin is saml", func() {
				It("adds the user", func() {
					u.addUser(SAMLOrigin, username)
					Expect(u.SamlUsers).To(HaveLen(1))
					Expect(u.SamlUsers[0]).To(Equal(username))
					Expect(u.Users).To(HaveLen(0))
					Expect(u.LDAPUsers).To(HaveLen(0))
				})
			})

			When("the origin is ldap", func() {
				It("adds the user", func() {
					u.addUser(LDAPOrigin, username)
					Expect(u.LDAPUsers).To(HaveLen(1))
					Expect(u.LDAPUsers[0]).To(Equal(username))
					Expect(u.SamlUsers).To(HaveLen(0))
					Expect(u.Users).To(HaveLen(0))
				})
			})

			When("that user exists in the other two origins", func() {
				BeforeEach(func() {
					u.LDAPUsers = append(u.LDAPUsers, username)
					u.SamlUsers = append(u.SamlUsers, username)
				})

				It("adds the user", func() {
					u.addUser(InternalOrigin, username)
					Expect(u.Users).To(HaveLen(1))
					Expect(u.Users[0]).To(Equal(username))
					Expect(u.SamlUsers).To(HaveLen(1))
					Expect(u.LDAPUsers).To(HaveLen(1))
				})
			})
		})

		When("that user already exists", func() {
			BeforeEach(func() {
				// note that u.Users corresponds to InternalOrigin
				u.Users = append(u.Users, username)
			})

			It("should do nothing", func() {
				usersLen := len(u.Users)
				u.addUser(InternalOrigin, username)
				Expect(u.Users).To(HaveLen(usersLen))
			})
		})
	})

	Context("A UserMgmt is queried for a particular user", func() {
		When("the user exists", func() {
			When("in LDAP", func() {
				BeforeEach(func() {
					u.LDAPUsers = append(u.LDAPUsers, username)
				})

				It("returns true", func() {
					Expect(u.hasUser(LDAPOrigin, username)).To(BeTrue())
				})
			})

			When("in SAML", func() {
				BeforeEach(func() {
					u.SamlUsers = append(u.SamlUsers, username)
				})

				It("returns true", func() {
					Expect(u.hasUser(SAMLOrigin, username)).To(BeTrue())
				})
			})

			When("in Internal", func() {
				BeforeEach(func() {
					u.Users = append(u.Users, username)
				})

				It("returns true", func() {
					Expect(u.hasUser(InternalOrigin, username)).To(BeTrue())
				})
			})
		})

		When("that user does not already exist", func() {
			It("returns false", func() {
				Expect(u.hasUser(InternalOrigin, username)).To(BeFalse())
			})

			When("the user exists on every other origin", func() {
				BeforeEach(func() {
					u.Users = append(u.Users, username)
					u.LDAPUsers = append(u.LDAPUsers, username)
				})

				It("returns false", func() {
					Expect(u.hasUser(SAMLOrigin, username)).To(BeFalse())
				})
			})
		})
	})
})
