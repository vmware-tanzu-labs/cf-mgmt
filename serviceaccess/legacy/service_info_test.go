package legacy_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	. "github.com/vmwarepivotallabs/cf-mgmt/serviceaccess/legacy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceInfo", func() {
	Context("GetPlan", func() {
		It("Should find service plan", func() {
			info := &ServiceInfo{}
			info.AddPlan("p-mysql", &resource.ServicePlan{Name: "small"})
			info.AddPlan("p-mysql", &resource.ServicePlan{Name: "large"})
			plans, err := info.GetPlans("p-mysql", []string{"small"})
			Expect(err).To(Not(HaveOccurred()))
			Expect(len(plans)).To(Equal(1))
			Expect(plans[0].Name).To(Equal("small"))
		})

		It("Should return all service plans", func() {
			info := &ServiceInfo{}
			info.AddPlan("p-mysql", &resource.ServicePlan{Name: "small"})
			info.AddPlan("p-mysql", &resource.ServicePlan{Name: "large"})
			plans, err := info.GetPlans("p-mysql", []string{"*"})
			Expect(err).To(Not(HaveOccurred()))
			Expect(len(plans)).To(Equal(2))
		})

		It("Should return no service plans", func() {
			info := &ServiceInfo{}
			info.AddPlan("p-mysql", &resource.ServicePlan{Name: "small"})
			info.AddPlan("p-mysql", &resource.ServicePlan{Name: "large"})
			plans, err := info.GetPlans("p-mysql", []string{""})
			Expect(err).To(HaveOccurred())
			Expect(len(plans)).To(Equal(0))
		})

		It("Should error when plan doesn't exist", func() {
			info := &ServiceInfo{}
			info.AddPlan("p-mysql", &resource.ServicePlan{Name: "small"})
			info.AddPlan("p-mysql", &resource.ServicePlan{Name: "large"})
			_, err := info.GetPlans("p-mysql", []string{"blah"})
			Expect(err).To(HaveOccurred())
		})
	})
})
