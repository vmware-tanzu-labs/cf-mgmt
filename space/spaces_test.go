package space_test

import (
	"errors"
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"

	"time"

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

	Context("FindSpace()", func() {
		It("should return an space", func() {
			spaces := []cfclient.Space{
				cfclient.Space{
					Name:             "testSpace",
					OrganizationGuid: "testOrgGUID",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			space, err := spaceManager.FindSpace("testOrg", "testSpace")
			Expect(err).Should(BeNil())
			Expect(space).ShouldNot(BeNil())
			Expect(space.Name).Should(Equal("testSpace"))
		})
		It("should return an error if space not found", func() {
			spaces := []cfclient.Space{
				{
					Name: "testSpace",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Expect(err).Should(HaveOccurred())
		})

		It("should return an error if unable to get OrgGUID", func() {
			fakeOrgMgr.GetOrgGUIDReturns("", fmt.Errorf("test"))
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Expect(err).Should(HaveOccurred())
		})
		It("should return an error if unable to get Spaces", func() {
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(nil, fmt.Errorf("test"))
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("CreateSpaces()", func() {
		BeforeEach(func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				config.SpaceConfig{
					Space: "space1",
				},
				config.SpaceConfig{
					Space: "space2",
				},
			}, nil)
		})
		It("should create 2 spaces", func() {
			spaces := []cfclient.Space{}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			Expect(spaceManager.CreateSpaces()).Should(Succeed())

			Expect(fakeClient.CreateSpaceCallCount()).Should(Equal(2))
			var spaceNames []string
			spaceRequest := fakeClient.CreateSpaceArgsForCall(0)
			Expect(spaceRequest.OrganizationGuid).Should(Equal("testOrgGUID"))
			spaceNames = append(spaceNames, spaceRequest.Name)
			spaceRequest = fakeClient.CreateSpaceArgsForCall(1)
			Expect(spaceRequest.OrganizationGuid).Should(Equal("testOrgGUID"))
			spaceNames = append(spaceNames, spaceRequest.Name)
			Expect(spaceNames).Should(ConsistOf([]string{"space1", "space2"}))
		})

		It("should create 1 space", func() {
			spaces := []cfclient.Space{
				{
					Name:             "space1",
					OrganizationGuid: "testOrgGUID",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)

			Expect(spaceManager.CreateSpaces()).Should(Succeed())
			Expect(fakeClient.CreateSpaceCallCount()).Should(Equal(1))
			spaceRequest := fakeClient.CreateSpaceArgsForCall(0)
			Expect(spaceRequest.OrganizationGuid).Should(Equal("testOrgGUID"))
			Expect(spaceRequest.Name).Should(Equal("space2"))
		})

		It("should create error if unable to get orgGUID", func() {
			fakeOrgMgr.GetOrgGUIDReturns("", errors.New("error1"))
			Expect(spaceManager.CreateSpaces()).ShouldNot(Succeed())
		})

		It("should rename a space", func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				config.SpaceConfig{
					Space:         "new-space1",
					OriginalSpace: "space1",
				},
			}, nil)
			spaces := []cfclient.Space{
				{
					Name:             "space1",
					Guid:             "space1-guid",
					OrganizationGuid: "testOrgGUID",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			Expect(spaceManager.CreateSpaces()).Should(Succeed())
			Expect(fakeClient.UpdateSpaceCallCount()).Should(Equal(1))
			spaceGUID, spaceRequest := fakeClient.UpdateSpaceArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space1-guid"))
			Expect(spaceRequest.Name).Should(Equal("new-space1"))
			Expect(spaceRequest.OrganizationGuid).Should(Equal("testOrgGUID"))
		})
	})

	Context("UpdateSpaces()", func() {
		BeforeEach(func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				config.SpaceConfig{
					Space:    "space1",
					AllowSSH: true,
				},
			}, nil)
		})
		It("should turn on allow ssh", func() {

			spaces := []cfclient.Space{
				{
					Name:             "space1",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space1GUID",
					AllowSSH:         false,
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			fakeClient.UpdateSpaceReturns(cfclient.Space{}, nil)

			err := spaceManager.UpdateSpaces()
			Expect(err).Should(BeNil())
			Expect(fakeClient.UpdateSpaceCallCount()).Should(Equal(1))
			spaceGUID, updateSpace := fakeClient.UpdateSpaceArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space1GUID"))
			Expect(updateSpace.OrganizationGuid).Should(Equal("testOrgGUID"))
			Expect(updateSpace.Name).Should(Equal("space1"))
			Expect(updateSpace.AllowSSH).Should(Equal(true))
		})

		It("should do nothing as ssh didn't change", func() {
			spaces := []cfclient.Space{
				{
					Name:             "space1",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space1GUID",
					AllowSSH:         true,
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			fakeClient.UpdateSpaceReturns(cfclient.Space{}, nil)

			err := spaceManager.UpdateSpaces()
			Expect(err).Should(BeNil())
			Expect(fakeClient.UpdateSpaceCallCount()).Should(Equal(0))
		})

		It("should turn on ssh temporarily", func() {
			future := time.Now().Add(time.Minute * 10)
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				config.SpaceConfig{
					Space:         "space1",
					AllowSSH:      false,
					AllowSSHUntil: future.Format(time.RFC3339),
				},
			}, nil)
			spaces := []cfclient.Space{
				{
					Name:             "space1",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space1GUID",
					AllowSSH:         false,
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			fakeClient.UpdateSpaceReturns(cfclient.Space{}, nil)

			err := spaceManager.UpdateSpaces()
			Expect(err).Should(BeNil())
			Expect(fakeClient.UpdateSpaceCallCount()).Should(Equal(1))
			_, spaceRequest := fakeClient.UpdateSpaceArgsForCall(0)
			Expect(spaceRequest.AllowSSH).To(BeTrue())
		})

		It("should turn off temporarily granted ssh", func() {
			past := time.Now().Add(time.Minute * -10)
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				config.SpaceConfig{
					Space:         "space1",
					AllowSSH:      false,
					AllowSSHUntil: past.Format(time.RFC3339),
				},
			}, nil)
			spaces := []cfclient.Space{
				{
					Name:             "space1",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space1GUID",
					AllowSSH:         true,
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			fakeClient.UpdateSpaceReturns(cfclient.Space{}, nil)

			err := spaceManager.UpdateSpaces()
			Expect(err).Should(BeNil())
			Expect(fakeClient.UpdateSpaceCallCount()).Should(Equal(1))
			_, spaceRequest := fakeClient.UpdateSpaceArgsForCall(0)
			Expect(spaceRequest.AllowSSH).To(BeFalse())
		})
		It("should do nothing as peek", func() {
			spaceManager.Peek = true
			spaces := []cfclient.Space{
				{
					Name:             "space1",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space1GUID",
					AllowSSH:         false,
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			fakeClient.UpdateSpaceReturns(cfclient.Space{}, nil)

			err := spaceManager.UpdateSpaces()
			Expect(err).Should(BeNil())
			Expect(fakeClient.UpdateSpaceCallCount()).Should(Equal(0))
		})

		It("should error on update space", func() {
			spaces := []cfclient.Space{
				{
					Name:             "space1",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space1GUID",
					AllowSSH:         false,
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			fakeClient.UpdateSpaceReturns(cfclient.Space{}, errors.New("error"))

			err := spaceManager.UpdateSpaces()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeClient.UpdateSpaceCallCount()).Should(Equal(1))
		})

	})

	Context("DeleteSpaces()", func() {
		BeforeEach(func() {
			fakeReader.SpacesReturns([]config.Spaces{
				config.Spaces{
					Spaces:             []string{"space1", "space2"},
					EnableDeleteSpaces: true,
				},
			}, nil)
			fakeReader.GetSpaceConfigReturns(&config.SpaceConfig{}, nil)
		})
		It("should delete 1", func() {
			spaces := []cfclient.Space{
				cfclient.Space{
					Name: "space1",
					Guid: "space1-guid",
				},
				cfclient.Space{
					Name: "space2",
					Guid: "space2-guid",
				},
				cfclient.Space{
					Name:             "space3",
					Guid:             "space3-guid",
					OrganizationGuid: "test2-org-guid",
				},
			}
			fakeOrgMgr.FindOrgReturns(cfclient.Org{
				Name: "test2",
				Guid: "test2-org-guid",
			}, nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			fakeClient.DeleteSpaceReturns(nil)
			Expect(spaceManager.DeleteSpaces()).Should(Succeed())
			Expect(fakeClient.DeleteSpaceCallCount()).Should(Equal(1))
			spaceGUID, recursive, async := fakeClient.DeleteSpaceArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space3-guid"))
			Expect(recursive).Should(Equal(true))
			Expect(async).Should(Equal(false))
		})

		It("should error", func() {
			spaces := []cfclient.Space{
				cfclient.Space{
					Name: "space1",
					Guid: "space1-guid",
				},
				cfclient.Space{
					Name: "space2",
					Guid: "space2-guid",
				},
				cfclient.Space{
					Name:             "space3",
					Guid:             "space3-guid",
					OrganizationGuid: "test2-org-guid",
				},
			}
			fakeOrgMgr.FindOrgReturns(cfclient.Org{
				Name: "test2",
				Guid: "test2-org-guid",
			}, nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			fakeClient.DeleteSpaceReturns(errors.New("error"))
			Expect(spaceManager.DeleteSpaces()).ShouldNot(Succeed())
			Expect(fakeClient.DeleteSpaceCallCount()).Should(Equal(1))
			spaceGUID, recursive, async := fakeClient.DeleteSpaceArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space3-guid"))
			Expect(recursive).Should(Equal(true))
			Expect(async).Should(Equal(false))
		})

		It("should peek", func() {
			spaceManager.Peek = true
			spaces := []cfclient.Space{
				cfclient.Space{
					Name: "space1",
					Guid: "space1-guid",
				},
				cfclient.Space{
					Name: "space2",
					Guid: "space2-guid",
				},
				cfclient.Space{
					Name: "space3",
					Guid: "space3-guid",
				},
			}
			fakeOrgMgr.FindOrgReturns(cfclient.Org{
				Name: "test2",
				Guid: "test2-org-guid",
			}, nil)
			fakeClient.ListSpacesReturns(spaces, nil)
			fakeClient.DeleteSpaceReturns(nil)
			Expect(spaceManager.DeleteSpaces()).Should(Succeed())
			Expect(fakeClient.DeleteSpaceCallCount()).Should(Equal(0))
		})
	})
})
