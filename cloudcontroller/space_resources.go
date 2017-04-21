package cloudcontroller

func (s *SpaceResources) GetNextURL() string {
	return s.NextURL
}

func NewSpaceResources() Pagination {
	return &SpaceResources{}
}

func (s *SpaceResources) AddInstances(temp Pagination) {
	if x, ok := temp.(*SpaceResources); ok {
		s.Spaces = append(s.Spaces, x.Spaces...)
	}
}
