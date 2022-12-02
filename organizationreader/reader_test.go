package organizationreader_test

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	. "github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
)

var _ = Describe("given OrgManager", func() {
	var (
		fakeClient *orgfakes.FakeCFOrganizationClient
		orgReader  DefaultReader
		fakeReader *configfakes.FakeReader
	)

	BeforeEach(func() {
		fakeClient = new(orgfakes.FakeCFOrganizationClient)
		fakeReader = new(configfakes.FakeReader)
		orgReader = DefaultReader{
			Cfg:    fakeReader,
			Client: fakeClient,
			Peek:   false,
		}
	})

	Context("FindOrg()", func() {
		It("should return an org", func() {
			orgs := []*resource.Organization{
				{
					Name: "test",
				},
				{
					Name: "test2",
				},
			}
			fakeClient.ListAllReturns(orgs, nil)
			org, err := orgReader.FindOrg("test")
			Ω(err).Should(BeNil())
			Ω(org).ShouldNot(BeNil())
			Ω(org.Name).Should(Equal("test"))
		})
	})
	It("should return an error for not found org", func() {
		var orgs []*resource.Organization
		fakeClient.ListAllReturns(orgs, nil)
		_, err := orgReader.FindOrg("test")
		Ω(err).ShouldNot(BeNil())
	})
	It("should return an error", func() {
		fakeClient.ListAllReturns(nil, fmt.Errorf("test"))
		_, err := orgReader.FindOrg("test")
		Ω(err).ShouldNot(BeNil())
	})

	Context("GetOrgGUID()", func() {
		It("should return an GUID", func() {
			orgs := []*resource.Organization{
				{
					Name: "test",
					GUID: "theGUID",
				},
			}
			fakeClient.ListAllReturns(orgs, nil)
			guid, err := orgReader.GetOrgGUID("test")
			Ω(err).Should(BeNil())
			Ω(guid).ShouldNot(BeNil())
			Ω(guid).Should(Equal("theGUID"))
		})
	})

	It("should return an error", func() {
		fakeClient.ListAllReturns(nil, fmt.Errorf("test"))
		guid, err := orgReader.GetOrgGUID("test")
		Ω(err).ShouldNot(BeNil())
		Ω(guid).Should(Equal(""))
	})
})
