package space_test

import (
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"

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

	Context("FindSpace()", func() {
		It("should return an space", func() {
			spaces := []*resource.Space{
				newCFSpace("testSpace-guid", "testSpace", "testOrgGUID"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			space, err := spaceManager.FindSpace("testOrg", "testSpace")
			Expect(err).Should(BeNil())
			Expect(space).ShouldNot(BeNil())
			Expect(space.Name).Should(Equal("testSpace"))
		})
		It("should return an error if space not found", func() {
			spaces := []*resource.Space{
				newCFSpace("testSpace-guid", "testSpace", "testOrgGUID"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
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
			fakeSpaceClient.ListAllReturns(nil, fmt.Errorf("test"))
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("CreateSpaces()", func() {
		BeforeEach(func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Space: "space1",
				},
				{
					Space: "space2",
				},
			}, nil)
		})
		It("should create 2 spaces", func() {
			spaces := []*resource.Space{}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			newSpace1 := newCFSpace("space1-guid", "space1", "test1-org-guid")
			newSpace2 := newCFSpace("space2-guid", "space2", "test1-org-guid")
			fakeSpaceClient.CreateReturnsOnCall(0, newSpace1, nil)
			fakeSpaceClient.CreateReturnsOnCall(0, newSpace2, nil)
			Expect(spaceManager.CreateSpaces()).Should(Succeed())

			Expect(fakeSpaceClient.CreateCallCount()).Should(Equal(2))
			var spaceNames []string
			_, spaceRequest := fakeSpaceClient.CreateArgsForCall(0)
			Expect(spaceRequest.Relationships.Organization.Data.GUID).Should(Equal("testOrgGUID"))
			spaceNames = append(spaceNames, spaceRequest.Name)
			_, spaceRequest = fakeSpaceClient.CreateArgsForCall(1)
			Expect(spaceRequest.Relationships.Organization.Data.GUID).Should(Equal("testOrgGUID"))
			spaceNames = append(spaceNames, spaceRequest.Name)
			Expect(spaceNames).Should(ConsistOf([]string{"space1", "space2"}))
		})

		It("should create 1 space", func() {
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test1-org-guid"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("test1-org-guid", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)

			Expect(spaceManager.CreateSpaces()).Should(Succeed())
			Expect(fakeSpaceClient.CreateCallCount()).Should(Equal(1))
			_, spaceRequest := fakeSpaceClient.CreateArgsForCall(0)
			Expect(spaceRequest.Relationships.Organization.Data.GUID).Should(Equal("test1-org-guid"))
			Expect(spaceRequest.Name).Should(Equal("space2"))
		})

		It("should create error if unable to get orgGUID", func() {
			fakeOrgMgr.GetOrgGUIDReturns("", errors.New("error1"))
			Expect(spaceManager.CreateSpaces()).ShouldNot(Succeed())
		})

		It("should rename a space", func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Space:         "new-space1",
					OriginalSpace: "space1",
				},
			}, nil)
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test1-org-guid"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("test1-org-guid", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			Expect(spaceManager.CreateSpaces()).Should(Succeed())
			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(1))
			_, spaceGUID, spaceRequest := fakeSpaceClient.UpdateArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space1-guid"))
			Expect(spaceRequest.Name).Should(Equal("new-space1"))
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
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test1-org-guid"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("test1-org-guid", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			fakeSpaceFeatureClient.IsSSHEnabledReturns(false, nil)
			fakeSpaceFeatureClient.EnableSSHReturns(nil)

			err := spaceManager.UpdateSpaces()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceFeatureClient.EnableSSHCallCount()).Should(Equal(1))
			_, spaceGUID, sshEnabled := fakeSpaceFeatureClient.EnableSSHArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space1-guid"))
			Expect(sshEnabled).Should(Equal(true))
		})

		It("should do nothing as ssh didn't change", func() {
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test1-org-guid"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("test1-org-guid", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			fakeSpaceFeatureClient.IsSSHEnabledReturns(true, nil)
			fakeSpaceFeatureClient.EnableSSHReturns(nil)

			err := spaceManager.UpdateSpaces()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceClient.UpdateCallCount()).Should(Equal(0))
		})

		It("should turn on ssh temporarily", func() {
			future := time.Now().Add(time.Minute * 10)
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Space:         "space1",
					AllowSSH:      false,
					AllowSSHUntil: future.Format(time.RFC3339),
				},
			}, nil)
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test1-org-guid"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("test1-org-guid", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			fakeSpaceFeatureClient.IsSSHEnabledReturns(false, nil)
			fakeSpaceFeatureClient.EnableSSHReturns(nil)

			err := spaceManager.UpdateSpaces()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceFeatureClient.EnableSSHCallCount()).Should(Equal(1))
			_, spaceGUID, sshEnabled := fakeSpaceFeatureClient.EnableSSHArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space1-guid"))
			Expect(sshEnabled).Should(Equal(true))
		})

		It("should turn off temporarily granted ssh", func() {
			past := time.Now().Add(time.Minute * -10)
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Space:         "space1",
					AllowSSH:      false,
					AllowSSHUntil: past.Format(time.RFC3339),
				},
			}, nil)
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test1-org-guid"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("test1-org-guid", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			fakeSpaceFeatureClient.IsSSHEnabledReturns(true, nil)
			fakeSpaceFeatureClient.EnableSSHReturns(nil)

			err := spaceManager.UpdateSpaces()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceFeatureClient.EnableSSHCallCount()).Should(Equal(1))
			_, spaceGUID, sshEnabled := fakeSpaceFeatureClient.EnableSSHArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space1-guid"))
			Expect(sshEnabled).Should(Equal(false))
		})
		It("should do nothing as peek", func() {
			spaceManager.Peek = true
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test1-org-guid"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("test1-org-guid", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			fakeSpaceFeatureClient.IsSSHEnabledReturns(false, nil)

			err := spaceManager.UpdateSpaces()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceFeatureClient.EnableSSHCallCount()).Should(Equal(0))
		})

		It("should error on update space", func() {
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test1-org-guid"),
			}
			fakeOrgMgr.GetOrgGUIDReturns("test1-org-guid", nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			fakeSpaceFeatureClient.IsSSHEnabledReturns(false, nil)
			fakeSpaceFeatureClient.EnableSSHReturns(errors.New("error"))

			err := spaceManager.UpdateSpaces()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeSpaceFeatureClient.EnableSSHCallCount()).Should(Equal(1))
		})

	})

	Context("DeleteSpaces()", func() {
		BeforeEach(func() {
			fakeReader.SpacesReturns([]config.Spaces{
				{
					Spaces:             []string{"space1", "space2"},
					EnableDeleteSpaces: true,
				},
			}, nil)
			fakeReader.GetSpaceConfigReturns(&config.SpaceConfig{}, nil)
		})
		It("should delete 1", func() {
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test1-org-guid"),
				newCFSpace("space2-guid", "space2", "test1-org-guid"),
				newCFSpace("space3-guid", "space3", "test2-org-guid"),
			}
			fakeOrgMgr.FindOrgReturns(&resource.Organization{
				Name: "test2",
				GUID: "test2-org-guid",
			}, nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			fakeSpaceClient.DeleteReturns("job-guid", nil)
			fakeJobClient.PollCompleteReturns(nil)

			Expect(spaceManager.DeleteSpaces()).Should(Succeed())
			Expect(fakeSpaceClient.DeleteCallCount()).Should(Equal(1))
			_, spaceGUID := fakeSpaceClient.DeleteArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space3-guid"))
			_, jobGUID, _ := fakeJobClient.PollCompleteArgsForCall(0)
			Expect(jobGUID).Should(Equal("job-guid"))
		})

		It("should error", func() {
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test1-org-guid"),
				newCFSpace("space2-guid", "space2", "test1-org-guid"),
				newCFSpace("space3-guid", "space3", "test2-org-guid"),
			}
			fakeOrgMgr.FindOrgReturns(&resource.Organization{
				Name: "test2",
				GUID: "test2-org-guid",
			}, nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			fakeSpaceClient.DeleteReturns("", errors.New("error"))
			Expect(spaceManager.DeleteSpaces()).ShouldNot(Succeed())
			Expect(fakeSpaceClient.DeleteCallCount()).Should(Equal(1))
			_, spaceGUID := fakeSpaceClient.DeleteArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space3-guid"))
			Expect(fakeJobClient.PollCompleteCallCount()).Should(Equal(0))
		})

		It("should peek", func() {
			spaceManager.Peek = true
			spaces := []*resource.Space{
				newCFSpace("space1-guid", "space1", "test2-org-guid"),
				newCFSpace("space2-guid", "space2", "test2-org-guid"),
				newCFSpace("space3-guid", "space3", "test2-org-guid"),
			}
			fakeOrgMgr.FindOrgReturns(&resource.Organization{
				Name: "test2",
				GUID: "test2-org-guid",
			}, nil)
			fakeSpaceClient.ListAllReturns(spaces, nil)
			Expect(spaceManager.DeleteSpaces()).Should(Succeed())
			Expect(fakeSpaceClient.DeleteCallCount()).Should(Equal(0))
		})
	})
})

func newCFSpace(guid, name, orgGUID string) *resource.Space {
	return &resource.Space{
		GUID: guid,
		Name: name,
		Relationships: &resource.SpaceRelationships{
			Organization: &resource.ToOneRelationship{
				Data: &resource.Relationship{
					GUID: orgGUID,
				},
			},
		},
	}
}
