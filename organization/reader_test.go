package organization_test

import (
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	configfakes "github.com/pivotalservices/cf-mgmt/config/fakes"
	. "github.com/pivotalservices/cf-mgmt/organization"
	orgfakes "github.com/pivotalservices/cf-mgmt/organization/fakes"
)

var _ = Describe("given OrgManager", func() {
	var (
		fakeClient *orgfakes.FakeCFClient
		orgReader  DefaultReader
		fakeReader *configfakes.FakeReader
	)

	BeforeEach(func() {
		fakeClient = new(orgfakes.FakeCFClient)
		fakeReader = new(configfakes.FakeReader)
		orgReader = DefaultReader{
			Cfg:    fakeReader,
			Client: fakeClient,
			Peek:   false,
		}
	})

	Context("FindOrg()", func() {
		It("should return an org", func() {
			orgs := []cfclient.Org{
				{
					Name: "test",
				},
				{
					Name: "test2",
				},
			}
			fakeClient.ListOrgsReturns(orgs, nil)
			org, err := orgReader.FindOrg("test")
			Ω(err).Should(BeNil())
			Ω(org).ShouldNot(BeNil())
			Ω(org.Name).Should(Equal("test"))
		})
	})
	It("should return an error for unfound org", func() {
		orgs := []cfclient.Org{}
		fakeClient.ListOrgsReturns(orgs, nil)
		_, err := orgReader.FindOrg("test")
		Ω(err).ShouldNot(BeNil())
	})
	It("should return an error", func() {
		fakeClient.ListOrgsReturns(nil, fmt.Errorf("test"))
		_, err := orgReader.FindOrg("test")
		Ω(err).ShouldNot(BeNil())
	})

	Context("GetOrgGUID()", func() {
		It("should return an GUID", func() {
			orgs := []cfclient.Org{
				{
					Name: "test",
					Guid: "theGUID",
				},
			}
			fakeClient.ListOrgsReturns(orgs, nil)
			guid, err := orgReader.GetOrgGUID("test")
			Ω(err).Should(BeNil())
			Ω(guid).ShouldNot(BeNil())
			Ω(guid).Should(Equal("theGUID"))
		})
	})

	It("should return an error", func() {
		fakeClient.ListOrgsReturns(nil, fmt.Errorf("test"))
		guid, err := orgReader.GetOrgGUID("test")
		Ω(err).ShouldNot(BeNil())
		Ω(guid).Should(Equal(""))
	})
})
