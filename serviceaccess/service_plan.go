package serviceaccess

import (
	"strings"

	"github.com/xchapter7x/lo"
)

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
	s.orgs[strings.ToLower(visibility.OrgGUID)] = visibility
}

func (s *ServicePlanInfo) RemoveOrg(orgGUID string) {
	delete(s.orgs, strings.ToLower(orgGUID))
}

func (s *ServicePlanInfo) OrgHasAccess(orgGUID string) bool {
	_, ok := s.orgs[strings.ToLower(orgGUID)]
	if ok {
		return true
	}
	lo.G.Debugf("OrgGUID %s is not in %+v", strings.ToLower(orgGUID), s.asKeys(s.orgs))
	return false
}

func (s *ServicePlanInfo) asKeys(input map[string]*Visibility) []string {
	keys := []string{}
	for key := range input {
		keys = append(keys, key)
	}
	return keys
}
