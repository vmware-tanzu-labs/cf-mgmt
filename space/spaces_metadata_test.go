package space_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"

	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
	uaafakes "github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
)

var _ = Describe("given SpaceManager", func() {
	var (
		fakeUaa                *uaafakes.FakeManager
		fakeOrgMgr             *orgfakes.FakeReader
		fakeSpaceClient        *spacefakes.FakeCFSpaceClient
		fakeSpaceFeatureClient *spacefakes.FakeCFSpaceFeatureClient
		fakeOrganizationClient *spacefakes.FakeCFOrganizationClient
		fakeJobClient          *spacefakes.FakeCFJobClient
		spaceManager           space.DefaultManager
		fakeReader             *configfakes.FakeReader
	)

	BeforeEach(func() {
		fakeUaa = new(uaafakes.FakeManager)
		fakeOrgMgr = new(orgfakes.FakeReader)
		fakeSpaceClient = new(spacefakes.FakeCFSpaceClient)
		fakeSpaceFeatureClient = new(spacefakes.FakeCFSpaceFeatureClient)
		fakeOrganizationClient = new(spacefakes.FakeCFOrganizationClient)
		fakeJobClient = new(spacefakes.FakeCFJobClient)
		fakeReader = new(configfakes.FakeReader)
		spaceManager = space.DefaultManager{
			Cfg:                fakeReader,
			SpaceClient:        fakeSpaceClient,
			SpaceFeatureClient: fakeSpaceFeatureClient,
			OrgClient:          fakeOrganizationClient,
			JobClient:          fakeJobClient,
			UAAMgr:             fakeUaa,
			OrgReader:          fakeOrgMgr,
			Peek:               false,
		}
	})

	Context("UpdateSpacesMetadata()", func() {
		It("should add metadata label for given space", func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Space: "testSpace",
					Org:   "testOrg",
					Metadata: &config.Metadata{
						Labels: map[string]string{
							"test-label": "test-value",
						},
					},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				MetadataPrefix: "foo.bar",
			}, nil)
			spaces := []*resource.Space{
				newCFSpace("test-space-guid", "testSpace", "testOrgGUID"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(1))
			_, spaceGUID, spaceUpdate := fakeSpaceClient.UpdateArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(spaceUpdate).ShouldNot(BeNil())
			Expect(spaceUpdate.Metadata).ShouldNot(BeNil())
			Expect(spaceUpdate.Metadata.Labels).Should(HaveKey("foo.bar/test-label"))
			Expect(*spaceUpdate.Metadata.Labels["foo.bar/test-label"]).Should(Equal("test-value"))
		})

		It("should add metadata annotation for given space", func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Space: "testSpace",
					Org:   "testOrg",
					Metadata: &config.Metadata{
						Annotations: map[string]string{
							"test-annotation": "test-value",
						},
					},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				MetadataPrefix: "foo.bar",
			}, nil)
			spaces := []*resource.Space{
				newCFSpace("test-space-guid", "testSpace", "testOrgGUID"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(1))
			_, spaceGUID, spaceUpdate := fakeSpaceClient.UpdateArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(spaceUpdate).ShouldNot(BeNil())
			Expect(spaceUpdate.Metadata).ShouldNot(BeNil())
			Expect(spaceUpdate.Metadata.Annotations).Should(HaveKey("foo.bar/test-annotation"))
			Expect(*spaceUpdate.Metadata.Annotations["foo.bar/test-annotation"]).Should(Equal("test-value"))
		})

		It("should remove metadata label for given space", func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Space: "testSpace",
					Org:   "testOrg",
					Metadata: &config.Metadata{
						Labels: map[string]string{
							"test-label": "",
						},
					},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				MetadataPrefix: "foo.bar",
			}, nil)
			spaces := []*resource.Space{
				newCFSpace("test-space-guid", "testSpace", "testOrgGUID"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(1))
			_, spaceGUID, spaceUpdate := fakeSpaceClient.UpdateArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(spaceUpdate).ShouldNot(BeNil())
			Expect(spaceUpdate.Metadata).ShouldNot(BeNil())
			Expect(spaceUpdate.Metadata.Labels).Should(HaveKey("foo.bar/test-label"))
			Expect(spaceUpdate.Metadata.Labels["foo.bar/test-label"]).Should(BeNil())
		})

		It("should remove metadata annotation for given space", func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Space: "testSpace",
					Org:   "testOrg",
					Metadata: &config.Metadata{
						Annotations: map[string]string{
							"test-annotation": "",
						},
					},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				MetadataPrefix: "foo.bar",
			}, nil)
			spaces := []*resource.Space{
				newCFSpace("test-space-guid", "testSpace", "testOrgGUID"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(1))
			_, spaceGUID, spaceUpdate := fakeSpaceClient.UpdateArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(spaceUpdate).ShouldNot(BeNil())
			Expect(spaceUpdate.Metadata).ShouldNot(BeNil())
			Expect(spaceUpdate.Metadata.Annotations).Should(HaveKey("foo.bar/test-annotation"))
			Expect(spaceUpdate.Metadata.Annotations["foo.bar/test-annotation"]).Should(BeNil())
		})

	})

	Context("ClearMetadata()", func() {
		It("should remove metadata from given space", func() {
			space := newCFSpace("test-space-guid", "testSpace", "testOrgGUID")
			space.Metadata = resource.NewMetadata().WithLabel("foo.bar", "test-label", "bar")
			err := spaceManager.ClearMetadata(space, "test-org")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(1))
			_, spaceGUID, spaceUpdate := fakeSpaceClient.UpdateArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(spaceUpdate).ShouldNot(BeNil())
			Expect(spaceUpdate.Metadata).ShouldNot(BeNil())
			Expect(spaceUpdate.Metadata.Labels).Should(HaveKey("foo.bar/test-label"))
			Expect(spaceUpdate.Metadata.Labels["foo.bar/test-label"]).Should(BeNil())
		})
	})
})
