package serviceaccess_test

import (
	. "github.com/vmwarepivotallabs/cf-mgmt/serviceaccess"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Context("GetPlan", func() {
		It("Should find service plan", func() {
			service := &Service{
				Name: "p-mysql",
			}
			service.AddPlan(&ServicePlanInfo{Name: "small"})
			service.AddPlan(&ServicePlanInfo{Name: "large"})
			plans, err := service.GetPlans([]string{"small"})
			Expect(err).To(Not(HaveOccurred()))
			Expect(len(plans)).To(Equal(1))
			Expect(plans[0].Name).To(Equal("small"))
		})
	})

	It("Should return all service plans", func() {
		service := &Service{
			Name: "p-mysql",
		}
		service.AddPlan(&ServicePlanInfo{Name: "small"})
		service.AddPlan(&ServicePlanInfo{Name: "large"})
		plans, err := service.GetPlans([]string{"*"})
		Expect(err).To(Not(HaveOccurred()))
		Expect(len(plans)).To(Equal(2))
	})
	It("Should return no service plans", func() {
		service := &Service{
			Name: "p-mysql",
		}
		service.AddPlan(&ServicePlanInfo{Name: "small"})
		service.AddPlan(&ServicePlanInfo{Name: "large"})
		plans, err := service.GetPlans([]string{""})
		Expect(err).To(HaveOccurred())
		Expect(len(plans)).To(Equal(0))
	})

	It("Should error when plan doesn't exist", func() {
		service := &Service{
			Name: "p-mysql",
		}
		service.AddPlan(&ServicePlanInfo{Name: "small"})
		service.AddPlan(&ServicePlanInfo{Name: "large"})
		_, err := service.GetPlans([]string{"blah"})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring("No plans for for service p-mysql"))
	})
})
