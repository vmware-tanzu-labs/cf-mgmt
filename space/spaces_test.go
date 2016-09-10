package space_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cf-mgmt/space"
)

var _ = Describe("given SpaceManager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewManager("test.com", "token", "uaacToken")
			Ω(manager).ShouldNot(BeNil())
		})
	})
})
