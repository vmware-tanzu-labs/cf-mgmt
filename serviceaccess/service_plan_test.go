package serviceaccess_test

import (
	. "github.com/vmwarepivotallabs/cf-mgmt/serviceaccess"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServicePlan", func() {
	var servicePlan *ServicePlanInfo
	Context("Adding org", func() {
		It("exists", func() {
			servicePlan = &ServicePlanInfo{}
			servicePlan.AddOrg(&Visibility{
				OrgGUID:         "org-guid-123-ABC",
				ServicePlanGUID: "service-plan-guid",
			})
			Expect(servicePlan.OrgHasAccess("org-guid-123-abc")).To(BeTrue())
		})

		It("exists", func() {
			servicePlan = &ServicePlanInfo{}
			servicePlan.AddOrg(&Visibility{
				OrgGUID:         "org-guid-123-ABC",
				ServicePlanGUID: "service-plan-guid",
			})
			Expect(servicePlan.OrgHasAccess("org-guid-123-ABC")).To(BeTrue())
		})

		It("doesn't exists", func() {
			servicePlan = &ServicePlanInfo{}
			servicePlan.AddOrg(&Visibility{
				OrgGUID:         "org-guid-123-efg",
				ServicePlanGUID: "service-plan-guid",
			})
			Expect(servicePlan.OrgHasAccess("org-guid-123-ABC")).To(BeFalse())
		})
	})
})
