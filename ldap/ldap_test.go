package ldap_test

import (
	"os"

	. "github.com/pivotalservices/cf-mgmt/ldap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ldap", func() {
	var ldapManager Manager
	Describe("given a GetUserIDs", func() {
		BeforeEach(func() {
			var host string

			if os.Getenv("LDAP_PORT_389_TCP_ADDR") == "" {
				host = "127.0.0.1"
			} else {
				host = os.Getenv("LDAP_PORT_389_TCP_ADDR")
			}
			ldapManager = &DefaultManager{
				LdapBindPassword: "password",
				Config: Config{
					BindDN:            "cn=admin,dc=pivotal,dc=org",
					UserSearchBase:    "ou=users,dc=pivotal,dc=org",
					UserNameAttribute: "uid",
					UserMailAttribute: "mail",
					GroupSearchBase:   "ou=groups,dc=pivotal,dc=org",
					GroupAttribute:    "member",
					LdapHost:          host,
					LdapPort:          389,
				},
				//LDAP_PORT_389_TCP_ADDR

			}
		})
		Context("when called with a valid group", func() {
			It("then it should return 3 users", func() {
				users, err := ldapManager.GetUserIDs("space_developers")
				Ω(err).Should(BeNil())
				Ω(len(users)).Should(Equal(3))
			})
		})
		XContext("when called with a valid group with special characters", func() {
			It("then it should return 3 users", func() {
				_, err := ldapManager.GetUserIDs("PCF One Org (Owner)")
				Ω(err).Should(BeNil())
				//Ω(len(users)).Should(Equal(3))
			})
		})
	})
})
