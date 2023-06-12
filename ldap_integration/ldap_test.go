package ldap_integration_test

import (
	"os"
	"strconv"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
	. "github.com/vmwarepivotallabs/cf-mgmt/ldap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ldap", func() {
	var ldapManager *Manager
	Describe("given a ldap manager", func() {

		BeforeEach(func() {
			var host string
			var port int
			var err error
			if os.Getenv("LDAP_PORT_389_TCP_ADDR") == "" {
				host = "127.0.0.1"
				port = 389
			} else {
				host = os.Getenv("LDAP_PORT_389_TCP_ADDR")
				port, _ = strconv.Atoi(os.Getenv("LDAP_PORT_389_TCP_PORT"))
			}
			ldapConfig := &config.LdapConfig{
				BindDN:            "cn=admin,dc=pivotal,dc=org",
				BindPassword:      "password",
				UserSearchBase:    "ou=users,dc=pivotal,dc=org",
				UserNameAttribute: "uid",
				UserMailAttribute: "mail",
				GroupSearchBase:   "dc=pivotal,dc=org",
				GroupAttribute:    "member",
				LdapHost:          host,
				LdapPort:          port,
			}
			ldapManager, err = NewManager(ldapConfig)
			Expect(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			ldapManager.Close()
		})
		Context("when cn with special characters", func() {
			It("then it should return 1 Entry", func() {
				entry, err := ldapManager.GetUserByDN(`cn=Washburn\2c Caleb,ou=users,dc=pivotal,dc=org`)
				Expect(err).Should(BeNil())
				Expect(entry).ShouldNot(BeNil())
			})
			It("then it should return 1 Entry", func() {
				entry, err := ldapManager.GetUserByDN("cn=Ekın Toğulmoç 88588,ou=users,dc=pivotal,dc=org")
				Expect(err).Should(BeNil())
				Expect(entry).ShouldNot(BeNil())
			})
		})
		Context("when cn has a period", func() {
			It("then it should return 1 Entry", func() {
				entry, err := ldapManager.GetUserByDN("cn=Caleb A. Washburn,ou=users,dc=pivotal,dc=org")
				Expect(err).Should(BeNil())
				Expect(entry).ShouldNot(BeNil())
			})
		})
		Context("when called with a valid group", func() {
			It("then it should return 5 users", func() {
				users, err := ldapManager.GetUserDNs("space_developers")
				Expect(err).Should(BeNil())
				Expect(len(users)).Should(Equal(6))
				Expect(users).To(ConsistOf([]string{
					"cn=cwashburn,ou=users,dc=pivotal,dc=org",
					"cn=Washburn\\2C Caleb,ou=users,dc=pivotal,dc=org",
					"cn=special\\2C (char) - username,ou=users,dc=pivotal,dc=org",
					"cn=Caleb A. Washburn,ou=users,dc=pivotal,dc=org",
					"cn=cwashburn1,ou=users,dc=pivotal,dc=org",
					"cn=Washburn\\2C Caleb\\2C cwashburn,ou=users,dc=pivotal,dc=org"}))
			})
		})
		Context("when called with a valid group with special characters", func() {
			It("then it should return 4 users", func() {
				users, err := ldapManager.GetUserDNs("special (char) group,name")
				Expect(err).Should(BeNil())
				Expect(len(users)).Should(Equal(4))
			})
		})
		Context("GetUser()", func() {
			It("then it should return 1 user", func() {
				user, err := ldapManager.GetUserByID("cwashburn")
				Expect(err).Should(BeNil())
				Expect(user).ShouldNot(BeNil())
				Expect(user.UserID).Should(Equal("cwashburn"))
				Expect(user.UserDN).Should(Equal("cn=cwashburn,ou=users,dc=pivotal,dc=org"))
				Expect(user.Email).Should(Equal("cwashburn+cfmt@testdomain.com"))
			})
		})

		Describe("given a ldap manager with userObjectClass", func() {
			BeforeEach(func() {
				var host string
				var port int
				var err error
				if os.Getenv("LDAP_PORT_389_TCP_ADDR") == "" {
					host = "127.0.0.1"
					port = 389
				} else {
					host = os.Getenv("LDAP_PORT_389_TCP_ADDR")
					port, _ = strconv.Atoi(os.Getenv("LDAP_PORT_389_TCP_PORT"))
				}
				ldapConfig := &config.LdapConfig{
					BindDN:            "cn=admin,dc=pivotal,dc=org",
					BindPassword:      "password",
					UserSearchBase:    "ou=users,dc=pivotal,dc=org",
					UserNameAttribute: "uid",
					UserMailAttribute: "mail",
					GroupSearchBase:   "dc=pivotal,dc=org",
					GroupAttribute:    "member",
					LdapHost:          host,
					LdapPort:          port,
					UserObjectClass:   "inetOrgPerson",
				}
				ldapManager, err = NewManager(ldapConfig)
				Expect(err).ShouldNot(HaveOccurred())

			})
			AfterEach(func() {
				ldapManager.Close()
			})
			Context("when cn with special characters", func() {
				It("then it should return 1 Entry", func() {
					entry, err := ldapManager.GetUserByDN(`cn=Washburn\2c Caleb,ou=users,dc=pivotal,dc=org`)
					Expect(err).Should(BeNil())
					Expect(entry).ShouldNot(BeNil())
				})
			})
			Context("GetUser()", func() {
				It("then it should return 1 user", func() {
					user, err := ldapManager.GetUserByID("cwashburn")
					Expect(err).Should(BeNil())
					Expect(user).ShouldNot(BeNil())
					Expect(user.UserID).Should(Equal("cwashburn"))
					Expect(user.UserDN).Should(Equal("cn=cwashburn,ou=users,dc=pivotal,dc=org"))
					Expect(user.Email).Should(Equal("cwashburn+cfmt@testdomain.com"))
				})
			})

			Context("GetLdapUser()", func() {
				It("then it should return 1 user", func() {
					data, _ := os.ReadFile("./fixtures/user1.txt")
					user, err := ldapManager.GetUserByDN(string(data))
					Expect(err).Should(BeNil())
					Expect(user).ShouldNot(BeNil())
					Expect(user.UserID).Should(Equal("cwashburn2"))
					Expect(user.Email).Should(Equal("cwashburn+cfmt2@testdomain.com"))
				})
			})
		})
	})
})
