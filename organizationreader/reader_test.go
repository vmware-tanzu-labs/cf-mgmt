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
		fakeOrgClient *orgfakes.FakeCFOrgClient
		orgReader     DefaultReader
		fakeReader    *configfakes.FakeReader
	)

	BeforeEach(func() {
		fakeOrgClient = new(orgfakes.FakeCFOrgClient)
		fakeReader = new(configfakes.FakeReader)
		orgReader = DefaultReader{
			Cfg:       fakeReader,
			OrgClient: fakeOrgClient,
			Peek:      false,
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
			fakeOrgClient.ListAllReturns(orgs, nil)
			org, err := orgReader.FindOrg("test")
			Ω(err).Should(BeNil())
			Ω(org).ShouldNot(BeNil())
			Ω(org.Name).Should(Equal("test"))
		})
	})
	It("should return an error for unfound org", func() {
		orgs := []*resource.Organization{}
		fakeOrgClient.ListAllReturns(orgs, nil)
		_, err := orgReader.FindOrg("test")
		Ω(err).ShouldNot(BeNil())
	})
	It("should return an error", func() {
		fakeOrgClient.ListAllReturns(nil, fmt.Errorf("test"))
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
			fakeOrgClient.ListAllReturns(orgs, nil)
			guid, err := orgReader.GetOrgGUID("test")
			Ω(err).Should(BeNil())
			Ω(guid).ShouldNot(BeNil())
			Ω(guid).Should(Equal("theGUID"))
		})
	})

	It("should return an error", func() {
		fakeOrgClient.ListAllReturns(nil, fmt.Errorf("test"))
		guid, err := orgReader.GetOrgGUID("test")
		Ω(err).ShouldNot(BeNil())
		Ω(guid).Should(Equal(""))
	})
})
