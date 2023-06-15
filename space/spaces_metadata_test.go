package space_test

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
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
		fakeUaa      *uaafakes.FakeManager
		fakeOrgMgr   *orgfakes.FakeReader
		fakeClient   *spacefakes.FakeCFClient
		spaceManager space.DefaultManager
		fakeReader   *configfakes.FakeReader
	)

	BeforeEach(func() {
		fakeUaa = new(uaafakes.FakeManager)
		fakeOrgMgr = new(orgfakes.FakeReader)
		fakeClient = new(spacefakes.FakeCFClient)
		fakeReader = new(configfakes.FakeReader)
		spaceManager = space.DefaultManager{
			Cfg:       fakeReader,
			Client:    fakeClient,
			UAAMgr:    fakeUaa,
			OrgReader: fakeOrgMgr,
			Peek:      false,
		}
	})

	Context("UpdateSpacesMetadata()", func() {
		It("should add metadata label for given space", func() {
			fakeClient.SupportsMetadataAPIReturns(true, nil)
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
			spaces := []cfclient.Space{
				{
					Name:             "testSpace",
					Guid:             "test-space-guid",
					OrganizationGuid: "testOrgGUID",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeClient.UpdateSpaceMetadataCallCount()).Should(Equal(1))
			spaceGUID, metadata := fakeClient.UpdateSpaceMetadataArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(metadata).ShouldNot(BeNil())
			Expect(metadata.Labels).Should(HaveKeyWithValue("foo.bar/test-label", "test-value"))
		})

		It("should add metadata annotation for given space", func() {
			fakeClient.SupportsMetadataAPIReturns(true, nil)
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
			spaces := []cfclient.Space{
				{
					Name:             "testSpace",
					Guid:             "test-space-guid",
					OrganizationGuid: "testOrgGUID",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeClient.UpdateSpaceMetadataCallCount()).Should(Equal(1))
			spaceGUID, metadata := fakeClient.UpdateSpaceMetadataArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(metadata).ShouldNot(BeNil())
			Expect(metadata.Annotations).Should(HaveKeyWithValue("foo.bar/test-annotation", "test-value"))
		})

		It("should remove metadata label for given space", func() {
			fakeClient.SupportsMetadataAPIReturns(true, nil)
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
			spaces := []cfclient.Space{
				{
					Name:             "testSpace",
					Guid:             "test-space-guid",
					OrganizationGuid: "testOrgGUID",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeClient.UpdateSpaceMetadataCallCount()).Should(Equal(1))
			spaceGUID, metadata := fakeClient.UpdateSpaceMetadataArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(metadata).ShouldNot(BeNil())
			Expect(metadata.Labels).Should(HaveKey("foo.bar/test-label"))
			Expect(metadata.Labels["foo.bar/test-label"]).Should(BeNil())
		})

		It("should remove metadata annotation for given space", func() {
			fakeClient.SupportsMetadataAPIReturns(true, nil)
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
			spaces := []cfclient.Space{
				{
					Name:             "testSpace",
					Guid:             "test-space-guid",
					OrganizationGuid: "testOrgGUID",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeClient.UpdateSpaceMetadataCallCount()).Should(Equal(1))
			spaceGUID, metadata := fakeClient.UpdateSpaceMetadataArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(metadata).ShouldNot(BeNil())
			Expect(metadata.Annotations).Should(HaveKey("foo.bar/test-annotation"))
			Expect(metadata.Annotations["foo.bar/test-annotation"]).Should(BeNil())
		})

	})

	Context("ClearMetadata()", func() {
		It("should remove metadata from given space", func() {
			fakeClient.SupportsMetadataAPIReturns(true, nil)
			space := cfclient.Space{
				Guid: "space-guid",
			}
			err := spaceManager.ClearMetadata(space, "test-org")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.RemoveSpaceMetadataCallCount()).Should(Equal(1))
			Expect(fakeClient.RemoveSpaceMetadataArgsForCall(0)).Should(Equal("space-guid"))
		})
	})
})
