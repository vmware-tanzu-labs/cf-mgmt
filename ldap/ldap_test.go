package ldap_test

import (
	"os"
	"strconv"

	. "github.com/pivotalservices/cf-mgmt/ldap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ldap", func() {
	var ldapManager Manager
	var config *Config
	Describe("given a GetUserIDs", func() {
		BeforeEach(func() {
			var host string
			var port int
			if os.Getenv("LDAP_PORT_389_TCP_ADDR") == "" {
				host = "127.0.0.1"
				port = 10389
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
		Context("when called with a valid group", func() {
			It("then it should return 4 users", func() {
				users, err := ldapManager.GetUserIDs(config, "space_developers")
				Ω(err).Should(BeNil())
				Ω(len(users)).Should(Equal(4))
			})
		})
		Context("when called with a valid group with special characters", func() {
			It("then it should return 4 users", func() {
				users, err := ldapManager.GetUserIDs(config, "special (char) group,name")
				Ω(err).Should(BeNil())
				Ω(len(users)).Should(Equal(4))
			})
		})
	})
})
