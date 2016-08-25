package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cf-mgmt/utils"
)

var _ = Describe("given utils manager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewDefaultManager()
			Ω(manager).ShouldNot(BeNil())
		})
	})
})
