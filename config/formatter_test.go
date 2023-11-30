package config_test

import (
	"time"

	. "github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/util"

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
		It("Should return formatted value for 1.2T", func() {
			val, err := StringToMegabytes("1200000")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(val).Should(Equal("1.2T"))
		})
	})

	Context("ByteSize", func() {
		It("Should return formatted value for less than 1GB", func() {
			val := ByteSize(util.GetIntPointer(1))
			Expect(val).Should(Equal("1M"))
		})
		It("Should return formatted value for 1GB", func() {
			val := ByteSize(util.GetIntPointer(1024))
			Expect(val).Should(Equal("1G"))
		})
		It("Should return formatted value for 1200GB", func() {
			val := ByteSize(util.GetIntPointer(1200000))
			Expect(val).Should(Equal("1.2T"))
		})

	})
	Context("ToMegabytes", func() {
		It("Should return int value for less than 1GB", func() {
			val, err := ToMegabytes("1M")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*val).Should(Equal(1))
		})
		It("Should return int value for more than 1GB", func() {
			val, err := ToMegabytes("1G")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*val).Should(Equal(1024))
		})
		It("Should return value for decimal values of measurement", func() {
			val, err := ToMegabytes("1.2T")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*val).Should(Equal(1200000))
		})
	})

	Context("ToInteger", func() {
		It("Should return int value of 1", func() {
			val, err := ToInteger("1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*val).Should(Equal(1))
		})
		It("Should return int value for 1024", func() {
			val, err := ToInteger("1024")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*val).Should(Equal(1024))
		})
		It("Should return int value for 1200000", func() {
			val, err := ToInteger("1200000")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(*val).Should(Equal(1200000))
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
