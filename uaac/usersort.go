package uaac

import (
	"sort"
	"strings"

	"github.com/xchapter7x/lo"
)

//DirAsc  Sort direction ascending
const DirAsc = "ASC"

//DirDesc  Sort direction descending
const DirDesc = "DESC"

//UserListSorter Type exposing a sort function for uaac.UserList
type UserListSorter struct {
	userList      *UserList
	sortDirection string
}

// NewUserListSorter --
func NewUserListSorter(users *UserList, sortDirection string) *UserListSorter {
	return &UserListSorter{
		userList:      users,
		sortDirection: sortDirection,
	}
}

//Sort Sorts the given userlist by Origin attribute asc
func (sorter *UserListSorter) Sort() *UserList {
	lo.G.Infof("Sort list size: %d", len(sorter.userList.Users))
	lo.G.Infof("Sorting user list based on Origin %s", sorter.sortDirection)
	sort.Sort(sorter)
	return sorter.userList
}

func (sorter *UserListSorter) Len() int { return len(sorter.userList.Users) }
func (sorter *UserListSorter) Less(i, j int) bool {
	user1 := sorter.userList.Users[i]
	user2 := sorter.userList.Users[j]
	if strings.Compare(user1.Origin, user2.Origin) == 0 {
		return user1.UserName < user2.UserName
	}
	if sorter.sortDirection == DirAsc {
		return user1.Origin < user2.Origin
	}
	return user1.Origin > user2.Origin
}
func (sorter *UserListSorter) Swap(i, j int) {
	sorter.userList.Users[i], sorter.userList.Users[j] = sorter.userList.Users[j], sorter.userList.Users[i]
}
