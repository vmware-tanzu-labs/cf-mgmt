package config_test

import (
	"time"

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

	Context("FutureTime", func() {
		It("Should return same time if no format is specified", func() {
			t := time.Now()
			future, err := FutureTime(t, "")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(t.Format(time.RFC3339)).Should(Equal(future))
		})

		It("Should add 1 day", func() {
			t := time.Now()
			expectedFuture := t.Add(time.Hour * 24 * 1)
			future, err := FutureTime(t, "1D")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(expectedFuture.Format(time.RFC3339)).Should(Equal(future))
		})

		It("Should add 5 hours", func() {
			t := time.Now()
			expectedFuture := t.Add(time.Hour * 5)
			future, err := FutureTime(t, "5H")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(expectedFuture.Format(time.RFC3339)).Should(Equal(future))
		})

		It("Should add 50 minutes", func() {
			t := time.Now()
			expectedFuture := t.Add(time.Minute * 50)
			future, err := FutureTime(t, "50M")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(expectedFuture.Format(time.RFC3339)).Should(Equal(future))
		})

		It("Should error with invalid time format", func() {
			t := time.Now()
			_, err := FutureTime(t, "5X")
			Expect(err).Should(HaveOccurred())
			Expect(err).Should(MatchError("Time to add must have format like D, M or H"))
		})
	})
})
