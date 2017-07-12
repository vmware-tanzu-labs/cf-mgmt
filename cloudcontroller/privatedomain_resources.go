package cloudcontroller

func (p *PrivateDomainResources) GetNextURL() string {
	return p.NextURL
}

func NewPrivateDomainResource() Pagination {
	return &PrivateDomainResources{}
}

func (p *PrivateDomainResources) AddInstances(temp Pagination) {
	if x, ok := temp.(*PrivateDomainResources); ok {
		p.PrivateDomains = append(p.PrivateDomains, x.PrivateDomains...)
	}
}
