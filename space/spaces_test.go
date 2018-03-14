package space_test

import (
	"errors"
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"

	orgfakes "github.com/pivotalservices/cf-mgmt/organization/fakes"
	"github.com/pivotalservices/cf-mgmt/space"
	spacefakes "github.com/pivotalservices/cf-mgmt/space/fakes"
	uaafakes "github.com/pivotalservices/cf-mgmt/uaa/fakes"
)

var _ = Describe("given SpaceManager", func() {
	var (
		fakeUaa      *uaafakes.FakeManager
		fakeOrgMgr   *orgfakes.FakeManager
		fakeClient   *spacefakes.FakeCFClient
		spaceManager space.DefaultManager
	)

	BeforeEach(func() {
		fakeUaa = new(uaafakes.FakeManager)
		fakeOrgMgr = new(orgfakes.FakeManager)
		fakeClient = new(spacefakes.FakeCFClient)
		spaceManager = space.DefaultManager{
			Cfg:    config.NewManager("./fixtures/config"),
			Client: fakeClient,
			UAAMgr: fakeUaa,
			OrgMgr: fakeOrgMgr,
			Peek:   false,
		}
	})

	Context("FindSpace()", func() {
		It("should return an space", func() {
			spaces := []cfclient.Space{
				{
					Name: "testSpace",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			space, err := spaceManager.FindSpace("testOrg", "testSpace")
			Ω(err).Should(BeNil())
			Ω(space).ShouldNot(BeNil())
			Ω(space.Name).Should(Equal("testSpace"))
		})
		It("should return an error if space not found", func() {
			spaces := []cfclient.Space{
				{
					Name: "testSpace",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error if unable to get OrgGUID", func() {
			fakeOrgMgr.GetOrgGUIDReturns("", fmt.Errorf("test"))
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())
		})
		It("should return an error if unable to get Spaces", func() {
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(nil, fmt.Errorf("test"))
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("CreateSpaces()", func() {
		BeforeEach(func() {
			spaceManager.Cfg = config.NewManager("./fixtures/config")
		})

		It("should create 2 spaces", func() {
			spaces := []cfclient.Space{}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			Ω(spaceManager.CreateSpaces()).Should(Succeed())

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
					Name: "space1",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)

			Ω(spaceManager.CreateSpaces()).Should(Succeed())
			Expect(fakeClient.CreateSpaceCallCount()).Should(Equal(1))
			spaceRequest := fakeClient.CreateSpaceArgsForCall(0)
			Expect(spaceRequest.OrganizationGuid).Should(Equal("testOrgGUID"))
			Expect(spaceRequest.Name).Should(Equal("space2"))
		})

		It("should create error if unable to get orgGUID", func() {
			fakeOrgMgr.GetOrgGUIDReturns("", errors.New("error1"))
			Ω(spaceManager.CreateSpaces()).ShouldNot(Succeed())
		})
	})

	/*Context("CreateApplicationSecurityGroups()", func() {
		It("should bind a named asg", func() {
			spaceManager.Cfg = config.NewManager("./fixtures/asg-config")

			spaces := []cfclient.Space{
				{
					Name: "space1",
					Guid: "space1GUID",
				},
				{
					Name: "space2",
					Guid: "space2GUID",
				},
			}
			sgs := make(map[string]cfclient.SecGroup)
			sgs["test-asg"] = cfclient.SecGroup{Name: "test-asg", Guid: "SGGZZUID"}
			sgs["test-space1"] = cfclient.SecGroup{Name: "test-space1", Guid: "SGGUID"}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			fakeSecMgr.ListNonDefaultSecurityGroupsReturns(sgs, nil)

			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")
			Ω(err).Should(BeNil())
			Expect(fakeSecMgr.UpdateSecurityGroupCallCount()).Should(Equal(1))
			Expect(fakeSecMgr.AssignSecurityGroupToSpaceCallCount()).Should(Equal(2))
		})

		It("should create 1 asg", func() {
			spaces := []cfclient.Space{
				{
					Name: "space1",
					Guid: "space1GUID",
				},
				{
					Name: "space2",
					Guid: "space2GUID",
				},
			}
			sgs := make(map[string]cfclient.SecGroup)
			sgs["foo"] = cfclient.SecGroup{Name: "foo", Guid: "SG-FOO-GUID"}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			fakeSecMgr.ListNonDefaultSecurityGroupsReturns(sgs, nil)
			fakeSecMgr.CreateSecurityGroupReturns(&cfclient.SecGroup{}, nil)
			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")
			Ω(err).Should(BeNil())
			Expect(fakeSecMgr.CreateSecurityGroupCallCount()).Should(Equal(1))
			Expect(fakeSecMgr.AssignSecurityGroupToSpaceCallCount()).Should(Equal(1))
		})

		It("should create update 1 asg", func() {
			spaces := []cfclient.Space{
				{
					Name: "space1",
					Guid: "space1GUID",
				},
				{
					Name: "space2",
					Guid: "space2GUID",
				},
			}
			sgs := make(map[string]cfclient.SecGroup)
			sgs["test-space1"] = cfclient.SecGroup{Name: "test-space1", Guid: "SGGUID"}
			sgs["foo"] = cfclient.SecGroup{Name: "foo", Guid: "SG-FOO-GUID"}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			fakeSecMgr.ListNonDefaultSecurityGroupsReturns(sgs, nil)

			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")
			Ω(err).Should(BeNil())
			Expect(fakeSecMgr.UpdateSecurityGroupCallCount()).Should(Equal(1))
			Expect(fakeSecMgr.AssignSecurityGroupToSpaceCallCount()).Should(Equal(1))
		})
	})*/

	/*Context("CreateQuotas()", func() {
			It("should create 2 quotas", func() {
				spaces := []cfclient.Space{
					{
						Name:             "space1",
						OrganizationGuid: "testOrgGUID",
						Guid:             "space1GUID",
					},
					{
						Name:             "space2",
						OrganizationGuid: "testOrgGUID",
						Guid:             "space2GUID",
					},
				}
				quotas := make(map[string]string)
				mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
				mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
				mockCloudController.EXPECT().ListAllSpaceQuotasForOrg("testOrgGUID").Return(quotas, nil)
				mockCloudController.EXPECT().CreateSpaceQuota(cloudcontroller.SpaceQuotaEntity{
					OrgGUID: "testOrgGUID",
					QuotaEntity: cloudcontroller.QuotaEntity{
						Name:                    "space1",
						MemoryLimit:             10240,
						InstanceMemoryLimit:     -1,
						TotalRoutes:             10,
						TotalServices:           -1,
						PaidServicePlansAllowed: true,
						AppInstanceLimit:        -1,
						TotalReservedRoutePorts: 0,
						TotalPrivateDomains:     -1,
						TotalServiceKeys:        -1,
					}}).Return("space1QuotaGUID", nil)
				mockCloudController.EXPECT().AssignQuotaToSpace("space1GUID", "space1QuotaGUID")

				mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
				mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
				mockCloudController.EXPECT().ListAllSpaceQuotasForOrg("testOrgGUID").Return(quotas, nil)
				mockCloudController.EXPECT().CreateSpaceQuota(cloudcontroller.SpaceQuotaEntity{
					OrgGUID: "testOrgGUID",
					QuotaEntity: cloudcontroller.QuotaEntity{
						Name:                    "space2",
						MemoryLimit:             10240,
						InstanceMemoryLimit:     -1,
						TotalRoutes:             10,
						TotalServices:           -1,
						PaidServicePlansAllowed: true,
						AppInstanceLimit:        -1,
						TotalReservedRoutePorts: 0,
						TotalPrivateDomains:     -1,
						TotalServiceKeys:        -1,
					}}).Return("space2QuotaGUID", nil)
				mockCloudController.EXPECT().AssignQuotaToSpace("space2GUID", "space2QuotaGUID")
				err := spaceManager.CreateQuotas("./fixtures/config")
				Ω(err).Should(BeNil())
			})

			It("should update 2 quota", func() {
				spaces := []cfclient.Space{
					{
						Name:             "space1",
						OrganizationGuid: "testOrgGUID",
						Guid:             "space1GUID",
					},
					{
						Name:             "space2",
						OrganizationGuid: "testOrgGUID",
						Guid:             "space2GUID",
					},
				}
				quotas := make(map[string]string)
				quotas["space1"] = "space1QuotaGUID"
				quotas["space2"] = "space2QuotaGUID"
				mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
				mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
				mockCloudController.EXPECT().ListAllSpaceQuotasForOrg("testOrgGUID").Return(quotas, nil)
				mockCloudController.EXPECT().UpdateSpaceQuota("space1QuotaGUID", cloudcontroller.SpaceQuotaEntity{
					OrgGUID: "testOrgGUID",
					QuotaEntity: cloudcontroller.QuotaEntity{
						Name:                    "space1",
						MemoryLimit:             10240,
						InstanceMemoryLimit:     -1,
						TotalRoutes:             10,
						TotalServices:           -1,
						PaidServicePlansAllowed: true,
						AppInstanceLimit:        -1,
						TotalReservedRoutePorts: 0,
						TotalPrivateDomains:     -1,
						TotalServiceKeys:        -1,
					}}).Return(nil)
				mockCloudController.EXPECT().AssignQuotaToSpace("space1GUID", "space1QuotaGUID")

				mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
				mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
				mockCloudController.EXPECT().ListAllSpaceQuotasForOrg("testOrgGUID").Return(quotas, nil)
				mockCloudController.EXPECT().UpdateSpaceQuota("space2QuotaGUID", cloudcontroller.SpaceQuotaEntity{
					OrgGUID: "testOrgGUID",
					QuotaEntity: cloudcontroller.QuotaEntity{
						Name:                    "space2",
						MemoryLimit:             10240,
						InstanceMemoryLimit:     -1,
						TotalRoutes:             10,
						TotalServices:           -1,
						PaidServicePlansAllowed: true,
						AppInstanceLimit:        -1,
						TotalReservedRoutePorts: 0,
						TotalPrivateDomains:     -1,
						TotalServiceKeys:        -1,
					}}).Return(nil)
				mockCloudController.EXPECT().AssignQuotaToSpace("space2GUID", "space2QuotaGUID")
				err := spaceManager.CreateQuotas("./fixtures/config")
				Ω(err).Should(BeNil())
			})
		})

		Context("UpdateSpaces()", func() {
			spaces := []cfclient.Space{
				{
					Name:             "space1",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space1GUID",
				},
				{
					Name:             "space2",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space2GUID",
				},
			}
			It("should turn on allow ssh", func() {
				mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
				mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
				mockCloudController.EXPECT().UpdateSpaceSSH(true, "space1GUID").Return(nil)
				mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
				mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
				mockCloudController.EXPECT().UpdateSpaceSSH(true, "space2GUID").Return(nil)

				spaceManager.Cfg = config.NewManager("./fixtures/config")
				err := spaceManager.UpdateSpaces("./fixtures/config")
				Ω(err).Should(BeNil())
			})

			It("should not modify anything", func() {
				mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
				mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
				mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
				mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)

				spaceManager.Cfg = config.NewManager("./fixtures/config-sshoff")
				Ω(spaceManager.UpdateSpaces("./fixtures/config-sshoff")).Should(Succeed())
			})

			It("should error when UpdateSpaceSSH errors", func() {
				mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
				mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
				mockCloudController.EXPECT().UpdateSpaceSSH(true, "space1GUID").Return(fmt.Errorf("test"))

				Ω(spaceManager.UpdateSpaces("./fixtures/config")).ShouldNot(Succeed())
			})
		})

		Context("DeleteSpaces()", func() {
			BeforeEach(func() {
				spaceManager.Cfg = config.NewManager("./fixtures/config-delete")
			})

			It("should delete 1 and skip 1", func() {
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
				mockOrgMgr.EXPECT().FindOrg("test2").Return(cfclient.Org{
					Name: "test2",
					Guid: "test2-org-guid",
				}, nil)
				mockCloudController.EXPECT().ListSpaces("test2-org-guid").Return(spaces, nil)
				mockCloudController.EXPECT().DeleteSpace("space3-guid").Return(nil)
				Ω(spaceManager.DeleteSpaces("./fixtures/config-delete")).Should(Succeed())
			})
		})
	})
	*/
})
