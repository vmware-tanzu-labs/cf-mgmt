package isosegment

import (
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"
	. "github.com/pivotalservices/cf-mgmt/isosegment/test_data"
	mocks "github.com/pivotalservices/cf-mgmt/utils/mocks"
)

func newTestUpdater(m manager) *Updater {
	utilsMockManager := mocks.NewMockUtilsManager()
	PopulateWithTestData(utilsMockManager)
	return &Updater{
		Cfg:     config.NewManager("./fixtures/0001", utilsMockManager),
		DryRun:  false,
		CleanUp: true,
		cc:      m,
	}
}

var _ = Describe("Isolation Segments", func() {
	var (
		u    *Updater
		ctrl *gomock.Controller
		m    *Mockmanager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(testingT)
		m = NewMockmanager(ctrl)
		u = newTestUpdater(m)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Ensure() that segments exist", func() {
		Context("when there is an error retrieving isolation segments", func() {
			BeforeEach(func() {
				m.EXPECT().GetIsolationSegments().Return(nil, errors.New(""))
			})

			It("fails", func() {
				Ω(u.Ensure()).ShouldNot(Succeed())
			})
		})

		Context("when no segments exist", func() {
			BeforeEach(func() {
				m.EXPECT().GetIsolationSegments().Return(nil, nil)
			})

			It("creates two isolation segments", func() {
				m.EXPECT().CreateIsolationSegment("iso00")
				m.EXPECT().CreateIsolationSegment("iso01")
				Ω(u.Ensure()).Should(Succeed())
			})

			It("doesnt create isolation segments when DryRun is enabled", func() {
				u.DryRun = true
				m.EXPECT().CreateIsolationSegment(gomock.Any()).Times(0)
				Ω(u.Ensure()).Should(Succeed())
			})
		})

		Context("when extra segments exist and CleanUp is enabled", func() {
			It("deletes the unneeded isolation segments", func() {
				u.CleanUp = true
				seg := []Segment{{Name: "iso00"}, {Name: "iso01"}, {Name: "extra"}}
				m.EXPECT().GetIsolationSegments().Return(seg, nil)
				m.EXPECT().DeleteIsolationSegment("extra")
				Ω(u.Ensure()).Should(Succeed())
			})
		})

		Context("when extra segments exist and CleanUp is disabled", func() {
			It("does not delete the extra isolation segments", func() {
				u.CleanUp = false
				seg := []Segment{{Name: "iso00"}, {Name: "iso01"}, {Name: "extra"}}
				m.EXPECT().GetIsolationSegments().Return(seg, nil)
				m.EXPECT().DeleteIsolationSegment(gomock.Any()).Times(0)
				Ω(u.Ensure()).Should(Succeed())
			})
		})

		Context("when all segments exist", func() {
			BeforeEach(func() {
				seg := []Segment{{Name: "iso00"}, {Name: "iso01"}}
				m.EXPECT().GetIsolationSegments().Return(seg, nil)
			})

			It("creates no isolation segments", func() {
				m.EXPECT().CreateIsolationSegment(gomock.Any()).Times(0)
				Ω(u.Ensure()).Should(Succeed())
			})
		})

		Context("when some segments exist", func() {
			BeforeEach(func() {
				seg := []Segment{{Name: "iso00"}}
				m.EXPECT().GetIsolationSegments().Return(seg, nil)
			})

			It("creates no isolation segments", func() {
				m.EXPECT().CreateIsolationSegment("iso01")
				Ω(u.Ensure()).Should(Succeed())
			})
		})
	})

	Describe("Entitle() an org to isolation segments", func() {
		Context("when both orgs are already entitled to their isolation segments", func() {
			BeforeEach(func() {
				m.EXPECT().EntitledIsolationSegments("org2").Return([]Segment{{Name: "iso00"}}, nil)
				m.EXPECT().EntitledIsolationSegments("org1").Return([]Segment{{Name: "iso01"}}, nil)
			})

			It("makes no changes", func() {
				m.EXPECT().EnableOrgIsolation(gomock.Any(), gomock.Any()).Times(0)
				m.EXPECT().RevokeOrgIsolation(gomock.Any(), gomock.Any()).Times(0)
				Ω(u.Entitle()).Should(Succeed())
			})
		})

		Context("when no orgs have been entitled to their isolation segments", func() {
			BeforeEach(func() {
				m.EXPECT().EntitledIsolationSegments(gomock.Any()).Return(nil, nil).AnyTimes()
			})

			It("entitles both orgs to their isolation segments", func() {
				By("entitling org2 to iso00 (it's default segment)")
				m.EXPECT().EnableOrgIsolation("org2", "iso00")

				By("entitling org1 to iso00 (used by one of its spaces)")
				m.EXPECT().EnableOrgIsolation("org1", "iso01")
				Ω(u.Entitle()).Should(Succeed())
			})

			It("makes no change when DryRun is enabled", func() {
				u.DryRun = true
				m.EXPECT().EnableOrgIsolation(gomock.Any(), gomock.Any()).Times(0)
				Ω(u.Entitle()).Should(Succeed())
			})
		})

		Context("when org1 is already entitled but org2 is not", func() {
			BeforeEach(func() {
				m.EXPECT().EntitledIsolationSegments("org2").Return(nil, nil)
				m.EXPECT().EntitledIsolationSegments("org1").Return([]Segment{{Name: "iso01"}}, nil)
			})

			It("entitles org2 to iso00", func() {
				m.EXPECT().EnableOrgIsolation("org2", "iso00")
				Ω(u.Entitle()).Should(Succeed())
			})
		})

		Context("when org2 is entitled to an extra isolation segment", func() {
			BeforeEach(func() {
				m.EXPECT().EntitledIsolationSegments("org2").Return([]Segment{{Name: "iso00"}, {Name: "extra"}}, nil)
				m.EXPECT().EntitledIsolationSegments("org1").Return([]Segment{{Name: "iso01"}}, nil)
			})

			It("revokes org2's access to the extra isolation segment when CleanUp is enabled", func() {
				m.EXPECT().RevokeOrgIsolation("org2", "extra")
				Ω(u.Entitle()).Should(Succeed())
			})

			It("does not revoke access when CleanUp is disabled", func() {
				u.CleanUp = false
				m.EXPECT().RevokeOrgIsolation("org2", "extra").Times(0)
				Ω(u.Entitle()).Should(Succeed())
			})

			It("makes no changes when DryRun is enabled", func() {
				u.DryRun = true
				m.EXPECT().EnableOrgIsolation(gomock.Any(), gomock.Any()).Times(0)
				m.EXPECT().RevokeOrgIsolation(gomock.Any(), gomock.Any()).Times(0)
				Ω(u.Entitle()).Should(Succeed())
			})
		})
	})

	Describe("UpdateOrgs() default isolation segment", func() {
		Context("when org2 is configured to use iso00 by default, and org1's config does not have a default", func() {
			It("sets isolation segments correctly", func() {
				m.EXPECT().SetOrgIsolationSegment("org2", Segment{Name: "iso00"})
				m.EXPECT().SetOrgIsolationSegment("org1", Segment{})
				Ω(u.UpdateOrgs()).Should(Succeed())
			})
		})

		Context("when DryRun is enabled", func() {
			BeforeEach(func() {
				u.DryRun = true
			})

			It("does not modify org isolation segments", func() {
				m.EXPECT().SetOrgIsolationSegment(gomock.Any(), gomock.Any()).Times(0)
				Ω(u.UpdateOrgs()).Should(Succeed())
			})
		})

		Context("when there is an error setting the default isolation segment", func() {
			BeforeEach(func() {
				m.EXPECT().SetOrgIsolationSegment(gomock.Any(), gomock.Any()).Return(errors.New("fail"))
			})

			It("fails", func() {
				Ω(u.UpdateOrgs()).ShouldNot(Succeed())
			})
		})
	})

	Describe("UpdateSpaces() isolation segments", func() {
		Context("when org1space2 is configured to use iso01", func() {
			It("sets isolation segments correctly", func() {
				m.EXPECT().SetSpaceIsolationSegment("org1", "org1space1", Segment{})
				m.EXPECT().SetSpaceIsolationSegment("org1", "org1space2", Segment{Name: "iso01"})
				m.EXPECT().SetSpaceIsolationSegment("org2", "org2space1", Segment{})
				Ω(u.UpdateSpaces()).Should(Succeed())
			})
		})

		Context("when DryRun is enabled", func() {
			BeforeEach(func() {
				u.DryRun = true
			})

			It("does not modify space isolation segments", func() {
				m.EXPECT().SetSpaceIsolationSegment(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				Ω(u.UpdateSpaces()).Should(Succeed())
			})
		})

		Context("when there is an error setting space isolation segment", func() {
			BeforeEach(func() {
				m.EXPECT().SetSpaceIsolationSegment(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("fail"))
			})

			It("fails", func() {
				Ω(u.UpdateSpaces()).ShouldNot(Succeed())
			})
		})
	})
})
