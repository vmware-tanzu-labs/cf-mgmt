package uaa

import "fmt"

func (u *UserList) GetNextURL(url string) string {
	if u.StartIndex+u.ItemsPerPage >= u.TotalResults {
		return ""
	}
	return fmt.Sprintf("%s&startIndex=%d", url, u.StartIndex+u.ItemsPerPage)
}

func NewUserListResources() Pagination {
	return &UserList{}
}

func (u *UserList) AddInstances(temp Pagination) {
	if x, ok := temp.(*UserList); ok {
		u.Users = append(u.Users, x.Users...)
	}
}
