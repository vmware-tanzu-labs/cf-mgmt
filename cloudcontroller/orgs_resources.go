package cloudcontroller

func (s *Orgs) GetNextURL() string {
	return s.NextURL
}

func NewOrgResources() Pagination {
	return &Orgs{}
}

func (s *Orgs) AddInstances(temp Pagination) {
	if x, ok := temp.(*Orgs); ok {
		s.Orgs = append(s.Orgs, x.Orgs...)
	}
}
