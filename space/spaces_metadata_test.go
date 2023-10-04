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
		spaceManager           space.DefaultManager
		fakeReader             *configfakes.FakeReader
		fakeSpaceClient        *spacefakes.FakeCFSpaceClient
		fakeSpaceFeatureClient *spacefakes.FakeCFSpaceFeatureClient
	)

	BeforeEach(func() {
		fakeUaa = new(uaafakes.FakeManager)
		fakeOrgMgr = new(orgfakes.FakeReader)
		fakeReader = new(configfakes.FakeReader)
		fakeSpaceClient = new(spacefakes.FakeCFSpaceClient)
		fakeSpaceFeatureClient = new(spacefakes.FakeCFSpaceFeatureClient)
		spaceManager = space.DefaultManager{
			Cfg:                fakeReader,
			UAAMgr:             fakeUaa,
			OrgReader:          fakeOrgMgr,
			Peek:               false,
			SpaceClient:        fakeSpaceClient,
			SpaceFeatureClient: fakeSpaceFeatureClient,
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
				{
					Name: "testSpace",
					GUID: "test-space-guid",
					Relationships: &resource.SpaceRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "testOrgGUID",
							},
						},
					},
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(1))
			_, spaceGUID, spaceUpdate := fakeSpaceClient.UpdateArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(spaceUpdate.Name).Should(Equal("testSpace"))
			value := spaceUpdate.Metadata.Labels["foo.bar/test-label"]
			Expect(*value).Should(Equal("test-value"))
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
				{
					Name: "testSpace",
					GUID: "test-space-guid",
					Relationships: &resource.SpaceRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "testOrgGUID",
							},
						},
					},
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(1))
			_, spaceGUID, spaceUpdate := fakeSpaceClient.UpdateArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			value := spaceUpdate.Metadata.Annotations["foo.bar/test-annotation"]
			Expect(*value).Should(Equal("test-value"))
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
				{
					Name: "testSpace",
					GUID: "test-space-guid",
					Relationships: &resource.SpaceRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "testOrgGUID",
							},
						},
					},
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(1))
			_, spaceGUID, spaceUpdate := fakeSpaceClient.UpdateArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
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
				{
					Name: "testSpace",
					GUID: "test-space-guid",
					Relationships: &resource.SpaceRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "testOrgGUID",
							},
						},
					},
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			err := spaceManager.UpdateSpacesMetadata()
			Expect(err).Should(BeNil())

			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(1))
			_, spaceGUID, spaceUpdate := fakeSpaceClient.UpdateArgsForCall(0)
			Expect(spaceGUID).Should(Equal("test-space-guid"))
			Expect(spaceUpdate.Metadata.Annotations).Should(HaveKey("foo.bar/test-annotation"))
			Expect(spaceUpdate.Metadata.Annotations["foo.bar/test-annotation"]).Should(BeNil())
		})

	})
})
