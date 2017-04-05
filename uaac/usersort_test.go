package uaac_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cf-mgmt/uaac"
)

var _ = Describe("Given a UserListSorter", func() {
	Describe("Create new sorter", func() {
		It("should return new sorter", func() {
			userList := &UserList{Users: make([]User, 1)}
			sorter := NewUserListSorter(userList, DirAsc)
			Ω(sorter).ShouldNot(BeNil())
		})
	})

	var (
		userList *UserList
	)

	BeforeEach(func() {
		users := make([]User, 0)
		user := User{ID: "ID1", UserName: "user1", Origin: "uaa"}
		users = append(users, user)
		user = User{ID: "ID2", UserName: "user2", Origin: "saml"}
		users = append(users, user)
		user = User{ID: "ID3", UserName: "user3", Origin: "ldap"}
		users = append(users, user)
		user = User{ID: "ID4", UserName: "user4", Origin: "ldap"}
		users = append(users, user)
		userList = &UserList{
			Users: users,
		}
	})

	Context("Should sort user list", func() {
		It("By Origin ASC, Username ASC", func() {

			sorter := NewUserListSorter(userList, DirAsc)
			sortedList := sorter.Sort()
			Ω(len(sortedList.Users)).Should(BeEquivalentTo(4))
			sortedUsers := sortedList.Users

			expectedUser := sortedUsers[0]
			Ω(expectedUser.Origin).Should(Equal("ldap"))
			Ω(expectedUser.UserName).Should(Equal("user3"))

			expectedUser = sortedUsers[1]
			Ω(expectedUser.Origin).Should(Equal("ldap"))
			Ω(expectedUser.UserName).Should(Equal("user4"))

			expectedUser = sortedUsers[2]
			Ω(expectedUser.Origin).Should(Equal("saml"))
			Ω(expectedUser.UserName).Should(Equal("user2"))
		})

		It("By Origin DESC, Username ASC", func() {

			sorter := NewUserListSorter(userList, DirDesc)
			sortedList := sorter.Sort()
			Ω(len(sortedList.Users)).Should(BeEquivalentTo(4))
			sortedUsers := sortedList.Users

			expectedUser := sortedUsers[0]
			Ω(expectedUser.Origin).Should(Equal("uaa"))
			Ω(expectedUser.UserName).Should(Equal("user1"))

			expectedUser = sortedUsers[1]
			Ω(expectedUser.Origin).Should(Equal("saml"))
			Ω(expectedUser.UserName).Should(Equal("user2"))

			expectedUser = sortedUsers[2]
			Ω(expectedUser.Origin).Should(Equal("ldap"))
			Ω(expectedUser.UserName).Should(Equal("user3"))
		})
	})
})
