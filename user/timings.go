package user

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/xchapter7x/lo"
)

func NewCFClientTimer(client CFClient) *CFClientTimer {
	return &CFClientTimer{
		client:  client,
		timings: make(map[string][]Timing),
	}
}

type Timing struct {
	ElapsedTime time.Duration
	Args        interface{}
}

type CFClientTimer struct {
	client  CFClient
	timings map[string][]Timing
}

func (t *CFClientTimer) LogResults() {
	for key, timings := range t.timings {
		lo.G.Infof("Method %s called %d times", key, len(timings))
	}

	file, err := os.Create("timings.csv")
	if err != nil {
		log.Fatal(err)
	}

	csvwriter := csv.NewWriter(file)
	for key, timings := range t.timings {
		for _, timing := range timings {
			csvwriter.Write([]string{key, timing.ElapsedTime.String(), fmt.Sprintf("%v", timing.Args)})
		}
	}
	csvwriter.Flush()
	defer file.Close()
}

func (t *CFClientTimer) addTiming(key string, timing Timing) {
	t.timings[key] = append(t.timings[key], timing)
}

func (t *CFClientTimer) ListSpacesByQuery(query url.Values) ([]cfclient.Space, error) {
	start := time.Now()
	result, err := t.client.ListSpacesByQuery(query)
	duration := time.Since(start)
	t.addTiming("ListSpacesByQuery", Timing{ElapsedTime: duration, Args: query})
	return result, err
}

func (t *CFClientTimer) DeleteUser(userGuid string) error {
	start := time.Now()
	err := t.client.DeleteUser(userGuid)
	duration := time.Since(start)
	t.addTiming("DeleteUser", Timing{ElapsedTime: duration, Args: userGuid})
	return err
}

func (t *CFClientTimer) DeleteV3Role(roleGUID string) error {
	start := time.Now()
	err := t.client.DeleteV3Role(roleGUID)
	duration := time.Since(start)
	t.addTiming("DeleteV3Role", Timing{ElapsedTime: duration, Args: roleGUID})
	return err
}

func (t *CFClientTimer) ListV3SpaceRolesByGUIDAndType(spaceGUID string, roleType string) ([]cfclient.V3User, error) {
	start := time.Now()
	result, err := t.client.ListV3SpaceRolesByGUIDAndType(spaceGUID, roleType)
	duration := time.Since(start)
	t.addTiming("ListV3SpaceRolesByGUIDAndType", Timing{ElapsedTime: duration, Args: spaceGUID + "_" + roleType})
	return result, err
}

func (t *CFClientTimer) ListV3OrganizationRolesByGUIDAndType(orgGUID string, roleType string) ([]cfclient.V3User, error) {
	start := time.Now()
	result, err := t.client.ListV3OrganizationRolesByGUIDAndType(orgGUID, roleType)
	duration := time.Since(start)
	t.addTiming("ListV3OrganizationRolesByGUIDAndType", Timing{ElapsedTime: duration, Args: orgGUID + "_" + roleType})
	return result, err
}

func (t *CFClientTimer) CreateV3OrganizationRole(orgGUID string, userGUID string, roleType string) (*cfclient.V3Role, error) {
	start := time.Now()
	result, err := t.client.CreateV3OrganizationRole(orgGUID, userGUID, roleType)
	duration := time.Since(start)
	t.addTiming("CreateV3OrganizationRole", Timing{ElapsedTime: duration, Args: orgGUID + "_" + userGUID + "_" + roleType})
	return result, err
}

func (t *CFClientTimer) CreateV3SpaceRole(spaceGUID string, userGUID string, roleType string) (*cfclient.V3Role, error) {
	start := time.Now()
	result, err := t.client.CreateV3SpaceRole(spaceGUID, userGUID, roleType)
	duration := time.Since(start)
	t.addTiming("CreateV3SpaceRole", Timing{ElapsedTime: duration, Args: spaceGUID + "_" + userGUID + "_" + roleType})
	return result, err
}

func (t *CFClientTimer) SupportsSpaceSupporterRole() (bool, error) {
	start := time.Now()
	result, err := t.client.SupportsSpaceSupporterRole()
	duration := time.Since(start)
	t.addTiming("SupportsSpaceSupporterRole", Timing{ElapsedTime: duration, Args: ""})
	return result, err
}

func (t *CFClientTimer) ListV3RolesByQuery(query url.Values) ([]cfclient.V3Role, error) {
	start := time.Now()
	result, err := t.client.ListV3RolesByQuery(query)
	duration := time.Since(start)
	t.addTiming("ListV3RolesByQuery", Timing{ElapsedTime: duration, Args: query})
	return result, err
}

func (t *CFClientTimer) ListV3UsersByQuery(query url.Values) ([]cfclient.V3User, error) {
	start := time.Now()
	result, err := t.client.ListV3UsersByQuery(query)
	duration := time.Since(start)
	t.addTiming("ListV3UsersByQuery", Timing{ElapsedTime: duration, Args: query})
	return result, err
}
