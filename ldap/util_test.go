package ldap_test

import (
	. "github.com/vmwarepivotallabs/cf-mgmt/ldap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
	Context("ParseUserCN", func() {
		It("Should return valid cn when cn has comma included", func() {
			cn, searchBase, err := ParseUserCN(`cn=Caleb\2c Washburn,ou=users,dc=pivotal,dc=org`)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cn).Should(BeEquivalentTo("cn=Caleb, Washburn"))
			Expect(searchBase).Should(BeEquivalentTo("ou=users,dc=pivotal,dc=org"))
		})
		It("Should return valid cn when cn has comma included", func() {
			cn, searchBase, err := ParseUserCN(`cn=Caleb, Washburn,ou=users,dc=pivotal,dc=org`)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cn).Should(BeEquivalentTo("cn=Caleb, Washburn"))
			Expect(searchBase).Should(BeEquivalentTo("ou=users,dc=pivotal,dc=org"))
		})
		It("Should return valid cn when cn has multi-byte character", func() {
			cn, searchBase, err := ParseUserCN("cn=Ekın,ou=users,dc=pivotal,dc=org")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cn).Should(BeEquivalentTo("cn=Ekın"))
			Expect(searchBase).Should(BeEquivalentTo("ou=users,dc=pivotal,dc=org"))
		})
		It("Should return first cn when there are multiples", func() {
			cn, searchBase, err := ParseUserCN("cn=Caleb Washubrn,cn=admin,ou=users,dc=pivotal,dc=org")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cn).Should(BeEquivalentTo("cn=Caleb Washubrn"))
			Expect(searchBase).Should(BeEquivalentTo("ou=users,dc=pivotal,dc=org"))
		})

		It("Should search base when other attribute types are found", func() {
			cn, searchBase, err := ParseUserCN("uid=AAAAAA,ou=BBBBBB,ou=CCCCCC,o=DDDDD,c=EEEEEE")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(cn).Should(BeEquivalentTo("uid=AAAAAA"))
			Expect(searchBase).Should(BeEquivalentTo("ou=BBBBBB,ou=CCCCCC,o=DDDDD,c=EEEEEE"))
		})

	})
})
