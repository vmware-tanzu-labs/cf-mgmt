package cloudcontroller

func (s *OrgSpaceUsers) GetNextURL() string {
	return s.NextURL
}

func NewOrgSpaceUsers() Pagination {
	return &OrgSpaceUsers{}
}

func (s *OrgSpaceUsers) AddInstances(temp Pagination) {
	if x, ok := temp.(*OrgSpaceUsers); ok {
		s.Users = append(s.Users, x.Users...)
	}
}
