package serviceaccess

type ServicePlanInfo struct {
	GUID        string
	Name        string
	ServiceName string
	Public      bool
	orgs        map[string]*Visibility
}

type Visibility struct {
	OrgGUID         string
	ServicePlanGUID string
}

func (s *ServicePlanInfo) ListVisibilities() []Visibility {
	var result []Visibility
	for _, visibility := range s.orgs {
		result = append(result, *visibility)
	}
	return result
}
func (s *ServicePlanInfo) AddOrg(visibility *Visibility) {
	if s.orgs == nil {
		s.orgs = make(map[string]*Visibility)
	}
	s.orgs[visibility.OrgGUID] = visibility
}

func (s *ServicePlanInfo) RemoveOrg(orgGUID string) {
	delete(s.orgs, orgGUID)
}

func (s *ServicePlanInfo) OrgHasAccess(orgGUID string) bool {
	_, ok := s.orgs[orgGUID]
	return ok
}
