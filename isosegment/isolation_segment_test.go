package isosegment_test

import (
	"errors"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/isosegment"
	"github.com/vmwarepivotallabs/cf-mgmt/isosegment/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
)

var _ = Describe("Isolation Segments", func() {
	var (
		u               *isosegment.Updater
		fakeIsoClient   *fakes.FakeCFIsolationSegmentClient
		fakeOrgClient   *fakes.FakeCFOrganizationClient
		fakeSpaceClient *fakes.FakeCFSpaceClient
		orgReader       *orgfakes.FakeReader
		spaceManager    *spacefakes.FakeManager
	)
	BeforeEach(func() {
		fakeIsoClient = new(fakes.FakeCFIsolationSegmentClient)
		fakeOrgClient = new(fakes.FakeCFOrganizationClient)
		fakeSpaceClient = new(fakes.FakeCFSpaceClient)
		orgReader = new(orgfakes.FakeReader)
		spaceManager = new(spacefakes.FakeManager)
		u = &isosegment.Updater{
			Cfg:              config.NewManager("./fixtures/0001"),
			OrgReader:        orgReader,
			SpaceManager:     spaceManager,
			Peek:             false,
			CleanUp:          true,
			IsoSegmentClient: fakeIsoClient,
			OrgClient:        fakeOrgClient,
			SpaceClient:      fakeSpaceClient,
		}
	})

	Describe("Ensure() that segments exist", func() {
		Context("when there is an error retrieving isolation segments", func() {
			BeforeEach(func() {
				fakeIsoClient.ListAllReturns(nil, errors.New(""))
			})

			It("fails", func() {
				Expect(u.Create()).ShouldNot(Succeed())
			})
		})

		Context("when no segments exist", func() {
			BeforeEach(func() {
				fakeIsoClient.ListAllReturns([]*resource.IsolationSegment{{Name: "shared"}}, nil)
			})

			It("creates isolation segments", func() {
				Expect(u.Create()).Should(Succeed())
				Expect(fakeIsoClient.CreateCallCount()).Should(Equal(2))
				var createdIsoSegments []string
				_, iso1 := fakeIsoClient.CreateArgsForCall(0)
				_, iso2 := fakeIsoClient.CreateArgsForCall(1)
				createdIsoSegments = append(createdIsoSegments, iso1.Name)
				createdIsoSegments = append(createdIsoSegments, iso2.Name)
				Expect(createdIsoSegments).Should(ConsistOf([]string{"iso01", "default_iso"}))
			})

			It("doesnt create isolation segments when DryRun is enabled", func() {
				u.Peek = true
				Ω(u.Create()).Should(Succeed())
				Expect(fakeIsoClient.CreateCallCount()).Should(Equal(0))
			})
		})

		Context("when extra segments exist and CleanUp is enabled", func() {
			It("deletes the unneeded isolation segments", func() {
				u.CleanUp = true
				seg := []*resource.IsolationSegment{
					{Name: "iso00", GUID: "iso00_guid"},
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "extra", GUID: "extra_guid"},
				}
				fakeIsoClient.ListAllReturns(seg, nil)
				Ω(u.Remove()).Should(Succeed())
				Expect(fakeIsoClient.DeleteCallCount()).Should(Equal(2))
				var deletedIsoSegments []string
				_, iso1 := fakeIsoClient.DeleteArgsForCall(0)
				_, iso2 := fakeIsoClient.DeleteArgsForCall(1)
				deletedIsoSegments = append(deletedIsoSegments, iso1)
				deletedIsoSegments = append(deletedIsoSegments, iso2)
				Expect(deletedIsoSegments).Should(ConsistOf([]string{"iso00_guid", "extra_guid"}))
			})

			It("doesn't delete the unneeded isolation segments", func() {
				u.CleanUp = true
				u.Peek = true
				seg := []*resource.IsolationSegment{
					{Name: "iso00", GUID: "iso00_guid"},
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "extra", GUID: "extra_guid"},
				}
				fakeIsoClient.ListAllReturns(seg, nil)
				Ω(u.Remove()).Should(Succeed())
				Expect(fakeIsoClient.DeleteCallCount()).Should(Equal(0))

			})
		})

		Context("when extra segments exist and CleanUp is disabled", func() {
			It("does not delete the extra isolation segments", func() {
				u.CleanUp = false
				seg := []*resource.IsolationSegment{
					{Name: "iso00"},
					{Name: "iso01"},
					{Name: "extra"},
				}
				fakeIsoClient.ListAllReturns(seg, nil)
				Ω(u.Remove()).Should(Succeed())
				Expect(fakeIsoClient.DeleteCallCount()).Should(Equal(0))
			})
		})

		Context("when all segments exist", func() {
			BeforeEach(func() {
				seg := []*resource.IsolationSegment{
					{Name: "iso00"},
					{Name: "iso01"},
					{Name: "default_iso"},
				}
				fakeIsoClient.ListAllReturns(seg, nil)
			})

			It("creates no isolation segments", func() {
				Ω(u.Create()).Should(Succeed())
				Expect(fakeIsoClient.CreateCallCount()).Should(Equal(0))
			})
		})

		Context("when some segments exist", func() {
			BeforeEach(func() {
				seg := []*resource.IsolationSegment{{Name: "iso00"}}
				fakeIsoClient.ListAllReturns(seg, nil)
			})

			It("creates isolation segments", func() {
				Ω(u.Create()).Should(Succeed())
				Expect(fakeIsoClient.CreateCallCount()).Should(Equal(2))
				var createdIsoSegments []string
				_, iso1 := fakeIsoClient.CreateArgsForCall(0)
				_, iso2 := fakeIsoClient.CreateArgsForCall(1)
				createdIsoSegments = append(createdIsoSegments, iso1.Name, iso2.Name)
				Expect(createdIsoSegments).Should(ConsistOf([]string{"iso01", "default_iso"}))
			})
		})
	})

	Describe("Entitle() an org to isolation segments", func() {
		Context("when both orgs are already entitled to their isolation segments", func() {
			BeforeEach(func() {
				fakeIsoClient.ListAllReturnsOnCall(0, []*resource.IsolationSegment{
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "default_iso", GUID: "default_iso_guid"},
				}, nil)
				orgReader.FindOrgReturns(&resource.Organization{GUID: "orgGUID"}, nil)
			})

			It("makes no changes", func() {
				fakeIsoClient.ListAllReturnsOnCall(1, []*resource.IsolationSegment{{Name: "iso01"}}, nil)
				Ω(u.Entitle()).Should(Succeed())
			})
		})

		Context("when no orgs have been entitled to their isolation segments", func() {
			BeforeEach(func() {
				fakeIsoClient.ListAllReturnsOnCall(0, []*resource.IsolationSegment{
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "default_iso", GUID: "default_iso_guid"},
				}, nil)
				fakeIsoClient.ListAllReturnsOnCall(1, []*resource.IsolationSegment{}, nil)
			})

			It("entitles both orgs to their isolation segments", func() {
				By("entitling org1 to iso00 (used by one of its spaces)")
				orgReader.FindOrgReturns(&resource.Organization{Name: "org1", GUID: "org1_guid"}, nil)
				Ω(u.Entitle()).Should(Succeed())
				Expect(fakeIsoClient.EntitleOrganizationCallCount()).Should(Equal(2))
				var isoSegmentGUIDs []string
				_, isolationSegmentGUID, orgGUID := fakeIsoClient.EntitleOrganizationArgsForCall(0)
				isoSegmentGUIDs = append(isoSegmentGUIDs, isolationSegmentGUID)
				Expect(orgGUID).Should(Equal("org1_guid"))
				_, isolationSegmentGUID, orgGUID = fakeIsoClient.EntitleOrganizationArgsForCall(1)
				isoSegmentGUIDs = append(isoSegmentGUIDs, isolationSegmentGUID)
				Expect(orgGUID).Should(Equal("org1_guid"))
				Expect(isoSegmentGUIDs).Should(ConsistOf([]string{"iso01_guid", "default_iso_guid"}))
			})

			It("makes no change when DryRun is enabled", func() {
				u.Peek = true
				orgReader.FindOrgReturns(&resource.Organization{Name: "org1", GUID: "org1_guid"}, nil)
				Ω(u.Entitle()).Should(Succeed())
				Expect(fakeIsoClient.EntitleOrganizationCallCount()).Should(Equal(0))
			})
		})

		Context("when org2 is entitled to an extra isolation segment", func() {
			BeforeEach(func() {
				fakeIsoClient.ListAllReturns([]*resource.IsolationSegment{
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "default_iso", GUID: "default_iso_guid"},
					{Name: "extra", GUID: "extra_guid"}}, nil)

				fakeIsoClient.ListAllReturns([]*resource.IsolationSegment{
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "default_iso", GUID: "default_iso_guid"},
					{Name: "extra", GUID: "extra_guid"}}, nil)
				orgReader.FindOrgReturns(&resource.Organization{Name: "org1", GUID: "org1_guid"}, nil)
			})

			It("revokes org2's access to the extra isolation segment when CleanUp is enabled", func() {
				Ω(u.Unentitle()).Should(Succeed())
				Expect(fakeIsoClient.RevokeOrganizationCallCount()).Should(Equal(1))
				_, isoGUID, orgGUID := fakeIsoClient.RevokeOrganizationArgsForCall(0)
				Expect(isoGUID).Should(Equal("extra_guid"))
				Expect(orgGUID).Should(Equal("org1_guid"))
			})

			It("does not revoke access when CleanUp is disabled", func() {
				u.CleanUp = false
				Ω(u.Entitle()).Should(Succeed())
				Expect(fakeIsoClient.RevokeOrganizationCallCount()).Should(Equal(0))
			})

			It("makes no changes when DryRun is enabled", func() {
				u.Peek = true
				Ω(u.Entitle()).Should(Succeed())
				Expect(fakeIsoClient.RevokeOrganizationCallCount()).Should(Equal(0))
			})
		})
	})

	Describe("UpdateOrgs() default isolation segment", func() {
		Context("when org1 is configured to use iso00 by default", func() {
			It("sets isolation segments correctly", func() {
				orgReader.FindOrgReturns(&resource.Organization{Name: "org1", GUID: "org1_guid"}, nil)
				fakeIsoClient.ListAllReturns([]*resource.IsolationSegment{
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "default_iso", GUID: "default_iso_guid"},
				}, nil)
				Ω(u.UpdateOrgs()).Should(Succeed())
				Expect(fakeOrgClient.AssignDefaultIsolationSegmentCallCount()).Should(Equal(1))
				_, orgGUID, isoSegmentGUID := fakeOrgClient.AssignDefaultIsolationSegmentArgsForCall(0)
				Expect(orgGUID).Should(Equal("org1_guid"))
				Expect(isoSegmentGUID).Should(Equal("default_iso_guid"))
			})
		})

		Context("when org1's config does not have a default", func() {
			BeforeEach(func() {
				u.Cfg = config.NewManager("./fixtures/0003")
			})
			It("sets isolation segments correctly", func() {
				orgReader.FindOrgReturns(&resource.Organization{Name: "org1", GUID: "org1_guid"}, nil)
				fakeIsoClient.ListAllReturns([]*resource.IsolationSegment{
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "default_iso", GUID: "default_iso_guid"},
				}, nil)
				fakeOrgClient.GetDefaultIsolationSegmentReturns("foo", nil)
				Ω(u.UpdateOrgs()).Should(Succeed())
				Expect(fakeOrgClient.AssignDefaultIsolationSegmentCallCount()).Should(Equal(1))
				_, orgGUID, isoGUID := fakeOrgClient.AssignDefaultIsolationSegmentArgsForCall(0)
				Expect(orgGUID).Should(Equal("org1_guid"))
				Expect(isoGUID).Should(Equal(""))
			})
			Context("when DryRun is enabled", func() {
				BeforeEach(func() {
					u.Peek = true
				})

				It("does not modify org isolation segments", func() {
					orgReader.FindOrgReturns(&resource.Organization{Name: "org1", GUID: "org1_guid"}, nil)
					Ω(u.UpdateOrgs()).Should(Succeed())
					Expect(fakeOrgClient.AssignDefaultIsolationSegmentCallCount()).Should(Equal(0))
				})
			})
		})

		Context("when DryRun is enabled", func() {
			BeforeEach(func() {
				u.Peek = true
			})

			It("does not modify org isolation segments", func() {
				orgReader.FindOrgReturns(&resource.Organization{Name: "org1", GUID: "org1_guid"}, nil)
				Ω(u.UpdateOrgs()).Should(Succeed())
				Expect(fakeOrgClient.AssignDefaultIsolationSegmentCallCount()).Should(Equal(0))
			})
		})

		Context("when there is an error setting the default isolation segment", func() {
			It("fails", func() {

				orgReader.FindOrgReturns(&resource.Organization{Name: "org1", GUID: "org1_guid"}, nil)
				fakeIsoClient.ListAllReturns([]*resource.IsolationSegment{
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "default_iso", GUID: "default_iso_guid"},
				}, nil)
				fakeOrgClient.AssignDefaultIsolationSegmentReturns(errors.New("error"))
				Ω(u.UpdateOrgs()).ShouldNot(Succeed())
			})
		})
	})

	Describe("UpdateSpaces() isolation segments", func() {
		Context("when org1space2 is configured to use iso01", func() {
			It("sets isolation segments correctly", func() {
				fakeIsoClient.ListAllReturns([]*resource.IsolationSegment{
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "default_iso", GUID: "default_iso_guid"},
				}, nil)
				spaceManager.FindSpaceReturns(&resource.Space{Name: "org1space2", GUID: "space_guid"}, nil)
				Ω(u.UpdateSpaces()).Should(Succeed())
				Expect(fakeSpaceClient.AssignIsolationSegmentCallCount()).Should(Equal(1))
				_, spaceGUID, isolationSegmentGUID := fakeSpaceClient.AssignIsolationSegmentArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(isolationSegmentGUID).Should(Equal("iso01_guid"))
			})
		})

		Context("when org1space2 is configured to use no isosegment", func() {
			BeforeEach(func() {
				u.Cfg = config.NewManager("./fixtures/0002")
			})
			It("sets isolation segments correctly", func() {
				fakeIsoClient.ListAllReturns([]*resource.IsolationSegment{
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "default_iso", GUID: "default_iso_guid"},
				}, nil)
				fakeSpaceClient.GetAssignedIsolationSegmentReturns("foo", nil)
				spaceManager.FindSpaceReturns(&resource.Space{Name: "org1space2", GUID: "space_guid"}, nil)
				Ω(u.UpdateSpaces()).Should(Succeed())
				Expect(fakeSpaceClient.AssignIsolationSegmentCallCount()).Should(Equal(1))
				_, spaceGUID, isoGUID := fakeSpaceClient.AssignIsolationSegmentArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(isoGUID).Should(Equal(""))
			})
			Context("when DryRun is enabled", func() {
				BeforeEach(func() {
					u.Peek = true
				})

				It("does not modify space isolation segments", func() {
					spaceManager.FindSpaceReturns(&resource.Space{Name: "org1space2", GUID: "space_guid"}, nil)
					Ω(u.UpdateSpaces()).Should(Succeed())
					Expect(fakeSpaceClient.AssignIsolationSegmentCallCount()).Should(Equal(0))
				})
			})
		})

		Context("when DryRun is enabled", func() {
			BeforeEach(func() {
				u.Peek = true
			})

			It("does not modify space isolation segments", func() {
				spaceManager.FindSpaceReturns(&resource.Space{Name: "org1space2", GUID: "space_guid"}, nil)
				Ω(u.UpdateSpaces()).Should(Succeed())
				Expect(fakeSpaceClient.AssignIsolationSegmentCallCount()).Should(Equal(0))
			})
		})

		Context("when there is an error setting space isolation segment", func() {

			It("fails", func() {
				fakeIsoClient.ListAllReturns([]*resource.IsolationSegment{
					{Name: "iso01", GUID: "iso01_guid"},
					{Name: "default_iso", GUID: "default_iso_guid"},
				}, nil)
				spaceManager.FindSpaceReturns(&resource.Space{Name: "org1space2", GUID: "space_guid"}, nil)
				fakeSpaceClient.AssignIsolationSegmentReturns(errors.New("error"))
				Ω(u.UpdateSpaces()).ShouldNot(Succeed())
			})
		})
	})
})
