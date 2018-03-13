package spaceusers_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cf-mgmt/spaceusers"
	"github.com/pivotalservices/cf-mgmt/spaceusers/fakes"
)

var _ = Describe("given UserSpaces", func() {
	var (
		userManager Manager
		client      *fakes.FakeCFClient
		userList    []cfclient.User
	)
	BeforeEach(func() {
		client = new(fakes.FakeCFClient)
	})
	Context("User Manager()", func() {
		BeforeEach(func() {
			userManager = NewManager(client, nil, nil, nil, nil, false)
			userList = []cfclient.User{
				cfclient.User{
					Username: "hello",
					Guid:     "world",
				},
				cfclient.User{
					Username: "hello2",
					Guid:     "world2",
				},
			}
		})

		Context("Success", func() {
			It("Should succeed on RemoveSpaceAuditorByUsername", func() {
				err := userManager.RemoveSpaceAuditorByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should succeed on RemoveSpaceDeveloperByUsername", func() {
				err := userManager.RemoveSpaceDeveloperByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceDeveloperByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should succeed on RemoveSpaceManagerByUsername", func() {
				err := userManager.RemoveSpaceManagerByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceManagerByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceManagerByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should succeed on ListSpaceAuditors", func() {
				client.ListSpaceAuditorsReturns(userList, nil)
				users, err := userManager.ListSpaceAuditors("foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(users)).Should(Equal(2))
				Expect(users).Should(HaveKeyWithValue("hello", "world"))
				Expect(users).Should(HaveKeyWithValue("hello2", "world2"))
				Expect(client.ListSpaceAuditorsCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceAuditorsArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})
			It("Should succeed on ListSpaceDevelopers", func() {
				client.ListSpaceDevelopersReturns(userList, nil)
				users, err := userManager.ListSpaceDevelopers("foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(users)).Should(Equal(2))
				Expect(users).Should(HaveKeyWithValue("hello", "world"))
				Expect(users).Should(HaveKeyWithValue("hello2", "world2"))
				Expect(client.ListSpaceDevelopersCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceDevelopersArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})
			It("Should succeed on ListSpaceManagers", func() {
				client.ListSpaceManagersReturns(userList, nil)
				users, err := userManager.ListSpaceManagers("foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(users)).Should(Equal(2))
				Expect(users).Should(HaveKeyWithValue("hello", "world"))
				Expect(users).Should(HaveKeyWithValue("hello2", "world2"))
				Expect(client.ListSpaceManagersCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceManagersArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})

			It("Should succeed on AssociateSpaceAuditorByUsername", func() {
				client.AssociateSpaceAuditorByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceAuditorByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.AssociateSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("spaceGUID"))
				Expect(userName).To(Equal("userName"))

				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should succeed on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceDeveloperByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.AssociateSpaceDeveloperByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("spaceGUID"))
				Expect(userName).To(Equal("userName"))

				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should succeed on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceManagerByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))

				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})
		})

		Context("Peek", func() {
			BeforeEach(func() {
				userManager = NewManager(client, nil, nil, nil, nil, false)
			})
			It("Should succeed on RemoveSpaceAuditorByUsername", func() {
				err := userManager.RemoveSpaceAuditorByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).To(Equal(1))
			})
			It("Should succeed on RemoveSpaceDeveloperByUsername", func() {
				err := userManager.RemoveSpaceDeveloperByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameCallCount()).To(Equal(1))
			})
			It("Should succeed on RemoveSpaceManagerByUsername", func() {
				err := userManager.RemoveSpaceManagerByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceManagerByUsernameCallCount()).To(Equal(1))
			})
			It("Should succeed on AssociateSpaceAuditorByUsername", func() {
				client.AssociateSpaceAuditorByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceAuditorByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should succeed on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceDeveloperByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should succeed on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceManagerByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
		})
		Context("Error", func() {
			It("Should error on RemoveSpaceAuditorByUsername", func() {
				client.RemoveSpaceAuditorByUsernameReturns(errors.New("error"))
				err := userManager.RemoveSpaceAuditorByUsername("foo", "bar")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should error on RemoveSpaceDeveloperByUsername", func() {
				client.RemoveSpaceDeveloperByUsernameReturns(errors.New("error"))
				err := userManager.RemoveSpaceDeveloperByUsername("foo", "bar")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceDeveloperByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should error on RemoveSpaceManagerByUsername", func() {
				client.RemoveSpaceManagerByUsernameReturns(errors.New("error"))
				err := userManager.RemoveSpaceManagerByUsername("foo", "bar")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceManagerByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceManagerByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should error on ListSpaceAuditors", func() {
				client.ListSpaceAuditorsReturns(nil, errors.New("error"))
				_, err := userManager.ListSpaceAuditors("foo")
				Expect(err).Should(HaveOccurred())
				Expect(client.ListSpaceAuditorsCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceAuditorsArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})
			It("Should error on ListSpaceDevelopers", func() {
				client.ListSpaceDevelopersReturns(nil, errors.New("error"))
				_, err := userManager.ListSpaceDevelopers("foo")
				Expect(err).Should(HaveOccurred())
				Expect(client.ListSpaceDevelopersCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceDevelopersArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})
			It("Should error on ListSpaceManagers", func() {
				client.ListSpaceManagersReturns(nil, errors.New("error"))
				_, err := userManager.ListSpaceManagers("foo")
				Expect(err).Should(HaveOccurred())
				Expect(client.ListSpaceManagersCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceManagersArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})
			It("Should error on AssociateSpaceAuditorByUsername", func() {
				client.AssociateSpaceAuditorByUsernameReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditorByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceAuditorByUsername", func() {
				client.AssociateSpaceAuditorByUsernameReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditorByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloperByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloperByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceManagerByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceManagerByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
		})
	})
})
