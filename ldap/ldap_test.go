package ldap_test

import (
	"io/ioutil"
	"os"
	"strconv"

	. "github.com/pivotalservices/cf-mgmt/ldap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ldap", func() {
	var ldapManager Manager
	var config *Config
	Describe("loading configuration", func() {
		Context("when there is valid ldap.yml", func() {
			It("then it should return a config", func() {
				config, err := NewManager().GetConfig("./fixtures/config", "test")
				Ω(err).Should(BeNil())
				Ω(config).ShouldNot(BeNil())
				Ω(config.Enabled).Should(BeTrue())
			})
		})
		Context("when there is invalid ldap.yml", func() {
			It("then it should return a config", func() {
				config, err := NewManager().GetConfig("./fixtures/blah", "test")
				Ω(err).Should(HaveOccurred())
				Ω(config).Should(BeNil())
			})
		})
	})
	Describe("given a ldap manager", func() {
		BeforeEach(func() {
			var host string
			var port int
			if os.Getenv("LDAP_PORT_389_TCP_ADDR") == "" {
				host = "127.0.0.1"
				port = 389
			} else {
				host = os.Getenv("LDAP_PORT_389_TCP_ADDR")
				port, _ = strconv.Atoi(os.Getenv("LDAP_PORT_389_TCP_PORT"))
			}
			ldapManager = &DefaultManager{}
			config = &Config{
				BindDN:            "cn=admin,dc=pivotal,dc=org",
				BindPassword:      "password",
				UserSearchBase:    "ou=users,dc=pivotal,dc=org",
				UserNameAttribute: "uid",
				UserMailAttribute: "mail",
				GroupSearchBase:   "ou=groups,dc=pivotal,dc=org",
				GroupAttribute:    "member",
				LdapHost:          host,
				LdapPort:          port,
			}
		})
		Context("when ldap is unreachable", func() {
			BeforeEach(func() { config.LdapHost = "unreachable-host" })
			It("then GetUserIDs should return an error", func() {
				_, err := ldapManager.GetUserIDs(config, "space_developers")
				Ω(err).ShouldNot(BeNil())
			})
			It("then GetUserIDs should return an error", func() {
				_, err := ldapManager.GetLdapUser(config, "cn=Washburn, Caleb,ou=users,dc=pivotal,dc=org", "ou=users,dc=pivotal,dc=org")
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("when bad password", func() {
			BeforeEach(func() { config.BindPassword = "foo" })
			It("then LdapConnection should return an error", func() {
				_, err := ldapManager.LdapConnection(config)
				Ω(err).ShouldNot(BeNil())
			})
		})
		Context("when bind user id has spaces", func() {
			BeforeEach(func() {
				config.BindDN = "cn=bind_account,ou=something with spaces,dc=pivotal,dc=org"
				config.BindPassword = "password"
			})
			It("then LdapConnection should not return an error", func() {
				_, err := ldapManager.LdapConnection(config)
				Ω(err).Should(BeNil())
			})
		})
		Context("when cn with special characters", func() {
			It("then it should return 1 Entry", func() {
				entry, err := ldapManager.GetLdapUser(config, "cn=Washburn, Caleb,ou=users,dc=pivotal,dc=org", "ou=users,dc=pivotal,dc=org")
				Ω(err).Should(BeNil())
				Ω(entry).ShouldNot(BeNil())
			})
		})
		Context("when called with a valid group", func() {
			It("then it should return 4 users", func() {
				users, err := ldapManager.GetUserIDs(config, "space_developers")
				Ω(err).Should(BeNil())
				Ω(len(users) > 3).Should(BeTrue())
			})
		})
		Context("when called with a valid group with special characters", func() {
			It("then it should return 4 users", func() {
				users, err := ldapManager.GetUserIDs(config, "special (char) group,name")
				Ω(err).Should(BeNil())
				Ω(len(users)).Should(Equal(4))
			})
		})
		Context("GetUser()", func() {
			It("then it should return 1 user", func() {
				user, err := ldapManager.GetUser(config, "cwashburn")
				Ω(err).Should(BeNil())
				Ω(user).ShouldNot(BeNil())
				Ω(user.UserID).Should(Equal("cwashburn"))
				Ω(user.UserDN).Should(Equal("cn=cwashburn,ou=users,dc=pivotal,dc=org"))
				Ω(user.Email).Should(Equal("cwashburn+cfmt@testdomain.com"))
			})
		})

		Describe("given a ldap manager with userObjectClass", func() {
			BeforeEach(func() {
				var host string
				var port int
				if os.Getenv("LDAP_PORT_389_TCP_ADDR") == "" {
					host = "127.0.0.1"
					port = 389
				} else {
					host = os.Getenv("LDAP_PORT_389_TCP_ADDR")
					port, _ = strconv.Atoi(os.Getenv("LDAP_PORT_389_TCP_PORT"))
				}
				ldapManager = &DefaultManager{}
				config = &Config{
					BindDN:            "cn=admin,dc=pivotal,dc=org",
					BindPassword:      "password",
					UserSearchBase:    "ou=users,dc=pivotal,dc=org",
					UserNameAttribute: "uid",
					UserMailAttribute: "mail",
					GroupSearchBase:   "ou=groups,dc=pivotal,dc=org",
					GroupAttribute:    "member",
					LdapHost:          host,
					LdapPort:          port,
					UserObjectClass:   "inetOrgPerson",
				}
			})
			Context("when cn with special characters", func() {
				It("then it should return 1 Entry", func() {
					entry, err := ldapManager.GetLdapUser(config, "cn=Washburn, Caleb,ou=users,dc=pivotal,dc=org", "ou=users,dc=pivotal,dc=org")
					Ω(err).Should(BeNil())
					Ω(entry).ShouldNot(BeNil())
				})
			})
			Context("GetUser()", func() {
				It("then it should return 1 user", func() {
					user, err := ldapManager.GetUser(config, "cwashburn")
					Ω(err).Should(BeNil())
					Ω(user).ShouldNot(BeNil())
					Ω(user.UserID).Should(Equal("cwashburn"))
					Ω(user.UserDN).Should(Equal("cn=cwashburn,ou=users,dc=pivotal,dc=org"))
					Ω(user.Email).Should(Equal("cwashburn+cfmt@testdomain.com"))
				})
			})
		})
		Context("GetLdapUser()", func() {
			It("then it should return 1 user", func() {
				data, _ := ioutil.ReadFile("./fixtures/user1.txt")
				user, err := ldapManager.GetLdapUser(config, string(data), "ou=users,dc=pivotal,dc=org")
				Ω(err).Should(BeNil())
				Ω(user).ShouldNot(BeNil())
				Ω(user.UserID).Should(Equal("cwashburn2"))
				Ω(user.Email).Should(Equal("cwashburn+cfmt2@testdomain.com"))
			})
		})
	})
})
