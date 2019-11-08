package ldap_test

import (
	. "github.com/pivotalservices/cf-mgmt/ldap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	Context("ParseUserCN", func() {
		It("Should return valid cn when cn has comma included", func() {
			cn, err := ParseUserCN(`cn=Caleb\2c Washburn,ou=users,dc=pivotal,dc=org`)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cn).Should(BeEquivalentTo("cn=Caleb, Washburn"))
		})
		It("Should return valid cn when cn has multi-byte character", func() {
			cn, err := ParseUserCN("cn=Ekın,ou=users,dc=pivotal,dc=org")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cn).Should(BeEquivalentTo("cn=Ekın"))
		})
		It("Should return first cn when there are multiples", func() {
			cn, err := ParseUserCN("cn=Caleb Washubrn,cn=admin,ou=users,dc=pivotal,dc=org")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cn).Should(BeEquivalentTo("cn=Caleb Washubrn"))
		})
	})
})
