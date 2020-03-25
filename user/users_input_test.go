package user_test

import (
	. "github.com/vmwarepivotallabs/cf-mgmt/user"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UsersInput", func() {
	Context("Ensure users are unique", func() {
		It("Should return unique user list", func() {
			input := UsersInput{
				Users: []string{"test", "test", "test2"},
			}
			Expect(input.UniqueUsers()).Should(ConsistOf([]string{"test", "test2"}))
		})
		It("Should return unique saml user list", func() {
			input := UsersInput{
				SamlUsers: []string{"test", "test", "test2"},
			}
			Expect(input.UniqueSamlUsers()).Should(ConsistOf([]string{"test", "test2"}))
		})
		It("Should return unique ldap user list", func() {
			input := UsersInput{
				LdapUsers: []string{"test", "test", "test2"},
			}
			Expect(input.UniqueLdapUsers()).Should(ConsistOf([]string{"test", "test2"}))
		})
		It("Should return unique ldap group list", func() {
			input := UsersInput{
				LdapGroupNames: []string{"test", "test", "test2"},
			}
			Expect(input.UniqueLdapGroupNames()).Should(ConsistOf([]string{"test", "test2"}))
		})
	})
})
