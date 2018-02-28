package cloudcontroller

func (s *IsoSegments) GetNextURL() string {
	return s.NextURL
}

func NewIsoSegmentResources() Pagination {
	return &IsoSegments{}
}

func (s *IsoSegments) AddInstances(temp Pagination) {
	if x, ok := temp.(*IsoSegments); ok {
		s.IsoSegments = append(s.IsoSegments, x.IsoSegments...)
	}
}
