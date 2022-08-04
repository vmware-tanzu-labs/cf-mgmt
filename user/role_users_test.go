package user_test

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	. "github.com/vmwarepivotallabs/cf-mgmt/user"
	"github.com/vmwarepivotallabs/cf-mgmt/user/fakes"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
	uaafakes "github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
)

var _ = Describe("RoleUsers", func() {
	var (
		userManager *DefaultManager
		client      *fakes.FakeCFClient
		ldapFake    *fakes.FakeLdapManager
		uaaFake     *uaafakes.FakeManager
		fakeReader  *configfakes.FakeReader
		spaceFake   *spacefakes.FakeManager
		orgFake     *orgfakes.FakeReader
	)
	BeforeEach(func() {
		client = new(fakes.FakeCFClient)
		ldapFake = new(fakes.FakeLdapManager)
		uaaFake = new(uaafakes.FakeManager)
		fakeReader = new(configfakes.FakeReader)
		spaceFake = new(spacefakes.FakeManager)
		orgFake = new(orgfakes.FakeReader)
		userManager = &DefaultManager{
			Client:     client,
			Cfg:        fakeReader,
			UAAMgr:     uaaFake,
			LdapMgr:    ldapFake,
			SpaceMgr:   spaceFake,
			OrgReader:  orgFake,
			Peek:       false,
			LdapConfig: &config.LdapConfig{Origin: "ldap"},
		}
		userList := []cfclient.V3User{
			{
				Username: "hello",
				GUID:     "world",
			},
			{
				Username: "hello2",
				GUID:     "world2",
			},
		}
		uaaUsers := &uaa.Users{}
		uaaUsers.Add(uaa.User{
			Username: "test",
			Origin:   "uaa",
			GUID:     "test-guid",
		})
		uaaUsers.Add(uaa.User{
			Username: "test-2",
			Origin:   "uaa",
			GUID:     "test2-guid",
		})
		uaaUsers.Add(uaa.User{
			Username: "hello",
			Origin:   "uaa",
			GUID:     "world",
		})
		uaaUsers.Add(uaa.User{
			Username: "hello2",
			Origin:   "uaa",
			GUID:     "world2",
		})
		userManager.UAAUsers = uaaUsers
		userMap := make(map[string]cfclient.V3User)
		for _, user := range userList {
			userMap[user.GUID] = user
		}
		userManager.CFUsers = userMap
	})
	Context("List Space Users", func() {
		BeforeEach(func() {

		})
		It("Return list of users by role", func() {
			Expect(true).Should(BeTrue())
		})
	})
})
