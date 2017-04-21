package cloudcontroller

func (s *SecurityGroupResources) GetNextURL() string {
	return s.NextURL
}

func NewSecurityGroupResources() Pagination {
	return &SecurityGroupResources{}
}

func (s *SecurityGroupResources) AddInstances(temp Pagination) {
	if x, ok := temp.(*SecurityGroupResources); ok {
		s.SecurityGroups = append(s.SecurityGroups, x.SecurityGroups...)
	}
}
