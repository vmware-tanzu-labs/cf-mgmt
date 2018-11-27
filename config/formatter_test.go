package config_test

import (
	. "github.com/pivotalservices/cf-mgmt/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Formatter", func() {
	Context("StringToMegabytes", func() {
		It("Should return formatted value for less than 1GB", func() {
			val, err := StringToMegabytes("1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).Should(Equal("1M"))
		})
		It("Should return formatted value for 1GB", func() {
			val, err := StringToMegabytes("1024")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).Should(Equal("1G"))
		})
	})

	Context("ByteSize", func() {
		It("Should return formatted value for less than 1GB", func() {
			val := ByteSize(1)
			Expect(val).Should(Equal("1M"))
		})
		It("Should return formatted value for 1GB", func() {
			val := ByteSize(1024)
			Expect(val).Should(Equal("1G"))
		})
	})
	Context("ToMegabytes", func() {
		It("Should return int value for less than 1GB", func() {
			val, err := ToMegabytes("1M")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).Should(Equal(1))
		})
		It("Should return int value for more than 1GB", func() {
			val, err := ToMegabytes("1G")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).Should(Equal(1024))
		})
	})
})
