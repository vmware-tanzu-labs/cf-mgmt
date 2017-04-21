package cloudcontroller

func (s *Quotas) GetNextURL() string {
	return s.NextURL
}

func NewQuotasResources() Pagination {
	return &Quotas{}
}

func (s *Quotas) AddInstances(temp Pagination) {
	if x, ok := temp.(*Quotas); ok {
		s.Quotas = append(s.Quotas, x.Quotas...)
	}
}
