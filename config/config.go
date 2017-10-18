// Package config provides utilities for reading and writing cf-mgmt's configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization/constants"
	"github.com/pivotalservices/cf-mgmt/space/constants"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

// DefaultProtectedOrgs lists the organizations that are considered protected
// and should never be deleted by cf-mgmt.
var DefaultProtectedOrgs = map[string]bool{
	"system":                  true,
	"p-spring-cloud-services": true,
	"splunk-nozzle-org":       true,
}

// Manager can read and write the cf-mgmt configuration.
type Manager interface {
	Updater
	Reader
}

// Updater is used to update the cf-mgmt configuration.
type Updater interface {
	AddOrgToConfig(orgConfig *OrgConfig) error
	AddSpaceToConfig(spaceConfig *SpaceConfig) error
	CreateConfigIfNotExists(uaaOrigin string) error
	DeleteConfigIfExists() error
	AddUserToSpaceConfig(userName, roleType, spaceName, orgName string, isLdapUser bool) error
	AddUserToOrgConfig(userName, roleType, orgName string, isLdapUser bool) error
	AddPrivateDomainToOrgConfig(orgName, privateDomainName string) error
	UpdateQuotasInOrgConfig(orgName string, enableQuota bool, parameters map[string]string) error
	UpdateQuotasInSpaceConfig(orgName, spaceName string, enableQuota bool, parameters map[string]string) error
	// -- Non Public Facing Functions
	setFieldIn(inputStruct interface{}, field, val string) error
}

// Reader is used to read the cf-mgmt configuration.
type Reader interface {
	Orgs() (Orgs, error)
	Spaces() ([]Spaces, error)

	GetOrgConfigs() ([]OrgConfig, error)
	GetAnOrgConfig(orgName string) (*OrgConfig, error)

	GetSpaceConfigs() ([]SpaceConfig, error)
	GetASpaceConfig(orgName, spaceName string, loadDefaults bool) (*SpaceConfig, error)
	GetSpaceDefaults() (*SpaceConfig, error)

	// -- Non Public Facing Functions
	getSpaceConfigsLoadDefaultOption(loadSpaceDefaults bool) ([]SpaceConfig, error)
}

// yamlManager is the default implementation of Manager.
// It is backed by a directory of YAML files.
type yamlManager struct {
	ConfigDir string
	UtilsMgr  utils.Manager
}

// NewManager creates a Manager that is backed by a set of YAML
// files in the specified configuration directory.
func NewManager(configDir string, utilsMgr utils.Manager) Manager {
	return &yamlManager{
		ConfigDir: configDir,
		UtilsMgr:  utilsMgr,
	}
}

// Orgs reads the config for all orgs.
func (m *yamlManager) Orgs() (Orgs, error) {
	configFile := filepath.Join(m.ConfigDir, "orgs.yml")
	lo.G.Info("Processing org file", configFile)
	input := Orgs{}
	if err := m.UtilsMgr.LoadFile(configFile, &input); err != nil {
		return Orgs{}, err
	}
	return input, nil
}

// GetOrgConfigs reads all orgs from the cf-mgmt configuration.
func (m *yamlManager) GetOrgConfigs() ([]OrgConfig, error) {
	fs := m.UtilsMgr
	files, err := fs.FindFiles(m.ConfigDir, "orgConfig.yml")
	if err != nil {
		return nil, err
	}

	result := make([]OrgConfig, len(files))
	for i, f := range files {
		result[i].AppInstanceLimit = -1
		result[i].TotalReservedRoutePorts = 0
		result[i].TotalPrivateDomains = -1
		result[i].TotalServiceKeys = -1
		if err = fs.LoadFile(f, &result[i]); err != nil {
			lo.G.Error(err)
			return nil, err
		}
	}

	return result, nil
}

func (m *yamlManager) Spaces() ([]Spaces, error) {
	fs := m.UtilsMgr
	files, err := fs.FindFiles(m.ConfigDir, "spaces.yml")
	if err != nil {
		return nil, err
	}

	spaceList := make([]Spaces, len(files))
	for i, f := range files {
		lo.G.Info("Processing space file", f)

		if err = fs.LoadFile(f, &spaceList[i]); err != nil {
			lo.G.Errorf("reading config for space %s: %v", f, err)
			return nil, err
		}
	}
	return spaceList, nil
}

func (m *yamlManager) getSpaceConfigsLoadDefaultOption(loadSpaceDefaults bool) ([]SpaceConfig, error) {
	fs := m.UtilsMgr

	spaceDefaults := SpaceConfig{}

	if loadSpaceDefaults {
		fs.LoadFile(filepath.Join(m.ConfigDir, "spaceDefaults.yml"), &spaceDefaults)
	}

	files, err := fs.FindFiles(m.ConfigDir, "spaceConfig.yml")
	if err != nil {
		return nil, err
	}
	result := make([]SpaceConfig, len(files))
	for i, f := range files {
		result[i].AppInstanceLimit = -1
		result[i].TotalReservedRoutePorts = 0
		result[i].TotalPrivateDomains = -1
		result[i].TotalServiceKeys = -1

		if err = fs.LoadFile(f, &result[i]); err != nil {
			return nil, err
		}

		if loadSpaceDefaults {
			result[i].Developer.LDAPUsers = append(result[i].Developer.LDAPUsers, spaceDefaults.Developer.LDAPUsers...)
			result[i].Developer.Users = append(result[i].Developer.Users, spaceDefaults.Developer.Users...)
			result[i].Developer.SamlUsers = append(result[i].Developer.SamlUsers, spaceDefaults.Developer.SamlUsers...)

			result[i].Auditor.LDAPUsers = append(result[i].Auditor.LDAPUsers, spaceDefaults.Auditor.LDAPUsers...)
			result[i].Auditor.Users = append(result[i].Auditor.Users, spaceDefaults.Auditor.Users...)
			result[i].Auditor.SamlUsers = append(result[i].Auditor.SamlUsers, spaceDefaults.Auditor.SamlUsers...)

			result[i].Manager.LDAPUsers = append(result[i].Manager.LDAPUsers, spaceDefaults.Manager.LDAPUsers...)
			result[i].Manager.Users = append(result[i].Manager.Users, spaceDefaults.Manager.Users...)
			result[i].Manager.SamlUsers = append(result[i].Manager.SamlUsers, spaceDefaults.Manager.SamlUsers...)

			result[i].Developer.LDAPGroups = append(result[i].GetDeveloperGroups(), spaceDefaults.GetDeveloperGroups()...)
			result[i].Auditor.LDAPGroups = append(result[i].GetAuditorGroups(), spaceDefaults.GetAuditorGroups()...)
			result[i].Manager.LDAPGroups = append(result[i].GetManagerGroups(), spaceDefaults.GetManagerGroups()...)
		}

		if result[i].EnableSecurityGroup {
			securityGroupFile := strings.Replace(f, "spaceConfig.yml", "security-group.json", -1)
			lo.G.Debug("Loading security group contents", securityGroupFile)
			bytes, err := m.UtilsMgr.LoadFileBytes(securityGroupFile)
			if err != nil {
				return nil, err
			}
			lo.G.Debug("setting security group contents", string(bytes))
			result[i].SecurityGroupContents = string(bytes)
		}
	}
	return result, nil
}

func (m *yamlManager) GetSpaceConfigs() ([]SpaceConfig, error) {
	return m.getSpaceConfigsLoadDefaultOption(true)
}

// GetSpaceDefaults returns the default space configuration, if one was provided.
// If no space defaults were configured, a nil config and a nil error are returned.
func (m *yamlManager) GetSpaceDefaults() (*SpaceConfig, error) {
	fp := filepath.Join(m.ConfigDir, "spaceDefaults.yml")
	fs := m.UtilsMgr

	if !fs.FileOrDirectoryExists(fp) {
		return nil, nil
	}
	result := SpaceConfig{}
	err := fs.LoadFile(fp, &result)
	return &result, err
}

// AddOrgToConfig adds an organization to the cf-mgmt configuration.
func (m *yamlManager) AddOrgToConfig(orgConfig *OrgConfig) error {
	orgFileName := filepath.Join(m.ConfigDir, "orgs.yml")
	orgName := orgConfig.Org
	if orgName == "" {
		return errors.New("cannot have an empty org name")
	}

	mgr := m.UtilsMgr
	orgList := &Orgs{}
	err := mgr.LoadFile(orgFileName, orgList)
	if err != nil {
		return err
	}

	if orgList.Contains(orgName) {
		lo.G.Infof("%s already added to config", orgName)
		return nil
	}
	lo.G.Infof("Adding org: %s ", orgName)
	orgList.Orgs = append(orgList.Orgs, orgName)
	if err = mgr.WriteFile(orgFileName, orgList); err != nil {
		return err
	}

	orgConfigFilePath := orgConfig.GetOrgConfigFilePath(m.ConfigDir, orgName)
	orgConfigFilenameAndPath := orgConfig.GetOrgConfigFilenameAndPath(m.ConfigDir, orgName)
	if err = os.MkdirAll(orgConfigFilePath, 0755); err != nil {
		return err
	}
	orgConfig.RemoveUsers = true
	orgConfig.RemovePrivateDomains = true
	mgr.WriteFile(orgConfigFilenameAndPath, orgConfig)
	return mgr.WriteFile(fmt.Sprintf("%s/spaces.yml", orgConfigFilePath), &Spaces{
		Org:                orgName,
		EnableDeleteSpaces: true,
	})
}

// AddSpaceToConfig adds a space to the cf-mgmt configuration, so long as a
// space with the specified name doesn't already exist.
func (m *yamlManager) AddSpaceToConfig(spaceConfig *SpaceConfig) error {
	orgName := spaceConfig.Org
	spaceFileName := filepath.Join(m.ConfigDir, orgName, "spaces.yml")
	spaceList := &Spaces{}
	spaceName := spaceConfig.Space
	mgr := m.UtilsMgr

	if err := mgr.LoadFile(spaceFileName, spaceList); err != nil {
		return err
	}
	if spaceList.Contains(spaceName) {
		lo.G.Infof("%s already added to config", spaceName)
		return nil
	}
	lo.G.Infof("Adding space: %s ", spaceName)
	spaceList.Spaces = append(spaceList.Spaces, spaceName)
	if err := mgr.WriteFile(spaceFileName, spaceList); err != nil {
		return err
	}
	spaceConfigPath := spaceConfig.GetSpaceConfigFilePath(m.ConfigDir, orgName, spaceName)
	if err := os.MkdirAll(spaceConfigPath, 0755); err != nil {
		return err
	}
	spaceConfig.RemoveUsers = true

	mgr.WriteFile(fmt.Sprintf("%s/spaceConfig.yml", spaceConfigPath), spaceConfig)
	mgr.WriteFileBytes(fmt.Sprintf("%s/security-group.json", spaceConfigPath), []byte("[]"))
	return nil
}

// CreateConfigIfNotExists initializes a new configuration directory.
// If the specified configuration directory already exists, it is left unmodified.
func (m *yamlManager) CreateConfigIfNotExists(uaaOrigin string) error {
	mgr := m.UtilsMgr
	if mgr.FileOrDirectoryExists(m.ConfigDir) {
		lo.G.Infof("Config directory %s already exists, skipping creation", m.ConfigDir)
		return nil
	}
	if err := os.MkdirAll(m.ConfigDir, 0755); err != nil {
		lo.G.Errorf("Error creating config directory %s. Error : %s", m.ConfigDir, err)
		return fmt.Errorf("cannot create directory %s: %v", m.ConfigDir, err)
	}
	lo.G.Infof("Config directory %s created", m.ConfigDir)
	mgr.WriteFile(fmt.Sprintf("%s/ldap.yml", m.ConfigDir), &ldap.Config{TLS: false, Origin: uaaOrigin})

	var protectedOrgs []string
	for protectedOrg := range DefaultProtectedOrgs {
		protectedOrgs = append(protectedOrgs, protectedOrg)
	}
	mgr.WriteFile(fmt.Sprintf("%s/orgs.yml", m.ConfigDir), &Orgs{
		EnableDeleteOrgs: true,
		ProtectedOrgs:    protectedOrgs,
	})
	mgr.WriteFile(fmt.Sprintf("%s/spaceDefaults.yml", m.ConfigDir), struct {
		Developer UserMgmt `yaml:"space-developer"`
		Manager   UserMgmt `yaml:"space-manager"`
		Auditor   UserMgmt `yaml:"space-auditor"`
	}{})
	return nil
}

// DeleteConfigIfExists deletes config directory if it exists.
func (m *yamlManager) DeleteConfigIfExists() error {
	utilsManager := m.UtilsMgr
	if !utilsManager.FileOrDirectoryExists(m.ConfigDir) {
		lo.G.Infof("%s doesn't exists, nothing to delete", m.ConfigDir)
		return nil
	}
	if err := os.RemoveAll(m.ConfigDir); err != nil {
		lo.G.Errorf("Error deleting config folder. Error: %s", err)
		return fmt.Errorf("cannot delete %s: %v", m.ConfigDir, err)
	}
	lo.G.Info("Config directory deleted")
	return nil
}

// GetASpaceConfig - Retrieves a single space config. loadDefaults is a boolean indicating if the spaceDefaults should
// be merged in together with the output spaceConfig.
func (m *yamlManager) GetASpaceConfig(orgName, spaceName string, loadDefaults bool) (*SpaceConfig, error) {
	// We would like to get the space config first.
	// This should not include the space defaults
	spaceConfigs, err := m.getSpaceConfigsLoadDefaultOption(loadDefaults)

	if err != nil {
		return nil, err
	}

	// Add our user to the space users
	// Find the space, in an Org, from the spaceConfigs array
	var targetSpaceConfig *SpaceConfig

	for _, spaceConfig := range spaceConfigs {
		if spaceConfig.Space == spaceName && spaceConfig.Org == orgName {
			targetSpaceConfig = &spaceConfig
			break
		}
	}

	// Check to ensure that our target space config was found.
	if targetSpaceConfig == nil {
		return nil, fmt.Errorf("The space %s was not found in %s", spaceName, orgName)
	}

	return targetSpaceConfig, nil
}

// AddUserToSpaceConfig adds a user to space in a given org.  isLdapUser specifies if the user is to be an ldap user
func (m *yamlManager) AddUserToSpaceConfig(userName, roleType, spaceName, orgName string, isLdapUser bool) error {
	// We would like to get the space config first.
	// This should not include the space defaults
	loadSpaceDefaults := false
	targetSpaceConfig, err := m.GetASpaceConfig(orgName, spaceName, loadSpaceDefaults)

	// Check to ensure that our target space config was found.
	if err != nil {
		return fmt.Errorf("Error retrieving space %s in Org %s because: %s", spaceName, orgName, err)
	} else if targetSpaceConfig == nil {
		return fmt.Errorf("The space %s was not found in Org %s", spaceName, orgName)
	}

	// Once we have the space, determine the user management role it fits into
	var userMgmtStruct *UserMgmt

	switch roleType {
	case space_constants.ROLE_SPACE_AUDITORS:
		userMgmtStruct = &targetSpaceConfig.Auditor
	case space_constants.ROLE_SPACE_DEVELOPERS:
		userMgmtStruct = &targetSpaceConfig.Developer
	case space_constants.ROLE_SPACE_MANAGERS:
		userMgmtStruct = &targetSpaceConfig.Manager
	default:
		return fmt.Errorf("Invalid Space Role: %s", roleType)
	}

	// Choose whether to use the user or the ldapUser to assign the role
	var targetUserRoleField *[]string
	if isLdapUser {
		targetUserRoleField = &userMgmtStruct.LDAPUsers
	} else {
		targetUserRoleField = &userMgmtStruct.Users
	}

	// Validate that the user is not already assigned with that role
	for _, user := range *targetUserRoleField {
		if user == userName {
			userType := ""
			if isLdapUser {
				userType = "LDAP "
			}
			return fmt.Errorf("%sUser %s already exists in %s/%s with the %s role", userType, userName, orgName, spaceName, roleType)
		}
	}

	// Add user into that role type
	*targetUserRoleField = append(*targetUserRoleField, userName)

	// Dump the file back out
	return m.UtilsMgr.WriteFile((*targetSpaceConfig).GetSpaceConfigFilenameAndPath(m.ConfigDir, orgName, spaceName), *targetSpaceConfig)
}

func (m *yamlManager) GetAnOrgConfig(orgName string) (*OrgConfig, error) {
	orgConfigs, err := m.GetOrgConfigs()
	if err != nil {
		return nil, err
	}

	// Add our user to the org users
	// Find the Org, from the orgConfigs array
	var targetOrgConfig *OrgConfig
	// Get a pointer to the target org
	for _, orgConfig := range orgConfigs {
		if orgConfig.Org == orgName {
			targetOrgConfig = &orgConfig
			break
		}
	}

	// Check to ensure that our target org config was found.
	if targetOrgConfig == nil {
		return nil, fmt.Errorf("The org %s was not found", orgName)
	}

	return targetOrgConfig, nil
}

// AddUserToOrgConfig adds a user to a given org.  isLdapUser specifies if the user is to be an ldap user
func (m *yamlManager) AddUserToOrgConfig(userName, roleType, orgName string, isLdapUser bool) error {
	orgConfig, err := m.GetAnOrgConfig(orgName)
	if err != nil {
		return err
	}

	// Once we have the org, determine the user management role it fits into
	var userMgmtStruct *UserMgmt

	switch roleType {
	case organization_constants.ROLE_ORG_AUDITORS:
		userMgmtStruct = &orgConfig.Auditor
	case organization_constants.ROLE_ORG_BILLING_MANAGERS:
		userMgmtStruct = &orgConfig.BillingManager
	case organization_constants.ROLE_ORG_MANAGERS:
		userMgmtStruct = &orgConfig.Manager
	default:
		return fmt.Errorf("Invalid Org Role: %s", roleType)
	}

	// Choose whether to use the user or the ldapUser to assign the role
	var targetUserRoleField *[]string
	if isLdapUser {
		targetUserRoleField = &userMgmtStruct.LDAPUsers
	} else {
		targetUserRoleField = &userMgmtStruct.Users
	}

	// Validate that the user is not already assigned with that role
	for _, user := range *targetUserRoleField {
		if user == userName {
			userType := ""
			if isLdapUser {
				userType = "LDAP "
			}
			return fmt.Errorf("%sUser %s already exists in %s with the %s role", userType, userName, orgName, roleType)
		}
	}

	// Add user into that role type
	*targetUserRoleField = append(*targetUserRoleField, userName)

	// Dump the file back out
	return m.UtilsMgr.WriteFile((*orgConfig).GetOrgConfigFilenameAndPath(m.ConfigDir, orgName), *orgConfig)
}

// AddOrgPrivateDomainToConfig adds a private domain to a given org.
func (m *yamlManager) AddPrivateDomainToOrgConfig(orgName, privateDomainName string) error {
	orgConfig, err := m.GetAnOrgConfig(orgName)
	if err != nil {
		return err
	}

	// Once we have the org, ensure the private domain doesn't already exist
	for _, storedPrivateDomainNamed := range orgConfig.PrivateDomains {
		if storedPrivateDomainNamed == privateDomainName {
			return fmt.Errorf("Private Domain Name %s already exists in %s", privateDomainName, orgName)
		}
	}

	// Add private domain to the org
	orgConfig.PrivateDomains = append(orgConfig.PrivateDomains, privateDomainName)

	// Dump the file back out
	return m.UtilsMgr.WriteFile((*orgConfig).GetOrgConfigFilenameAndPath(m.ConfigDir, orgName), *orgConfig)
}

// Helper Function - setFieldIn.
// Sets a field in an inputStruct specified by field, with a value specified by val
func (m *yamlManager) setFieldIn(inputStruct interface{}, field, val string) error {
	var err error
	ps := reflect.ValueOf(inputStruct)
	// struct
	s := ps.Elem()
	if s.Kind() == reflect.Struct {
		// exported field
		f := s.FieldByName(field)
		if f.IsValid() {
			// A Value can be changed only if it is
			// addressable and it is exportable (i.e. has Capital letter at the start)
			if f.CanSet() {
				switch f.Kind() {
				case reflect.Int:
					var outVal int64
					outVal, err = strconv.ParseInt(val, 10, 64)
					// If there's no error, set the field value
					if err == nil {
						f.SetInt(outVal)
					}
				case reflect.Bool:
					var outVal bool
					outVal, err = strconv.ParseBool(val)
					// If there's no error, set the field value
					if err == nil {
						f.SetBool(outVal)
					}
				default:
					err = fmt.Errorf("The parameter %s cannot be set - ensure it is a valid field", field)
				}

			} else {
				err = fmt.Errorf("The parameter %s cannot be set - ensure it is a valid field", field)
			}
		} else {
			err = fmt.Errorf("The parameter %s does not exist in the input structure", field)
		}
	} else {
		err = fmt.Errorf("Invalid structure input")
	}

	return err
}

// UpdateQuotasInOrgConfig updates the quotas specified in parameters and enables/disables the quotas
func (m *yamlManager) UpdateQuotasInOrgConfig(orgName string, enableQuota bool, parameters map[string]string) error {

	// Get the Org Cofig
	orgConfig, err := m.GetAnOrgConfig(orgName)
	if err != nil {
		return err
	}

	// Set the Enable Org Quota Field
	orgConfig.EnableOrgQuota = enableQuota

	// For each parameter key (field), find it in the OrgConfig struct
	// If the field is is found, update it with the parameter value.
	for field, value := range parameters {
		err := m.setFieldIn(orgConfig, field, value)
		// We should always just error out if any of the parameters fail to set
		if err != nil {
			return err
		}
	}

	// Dump the file back out
	return m.UtilsMgr.WriteFile((*orgConfig).GetOrgConfigFilenameAndPath(m.ConfigDir, orgName), *orgConfig)
}

// UpdateQuotasInSpaceConfig updates the quotas specified in parameters and enables/disables the quotas
func (m *yamlManager) UpdateQuotasInSpaceConfig(orgName, spaceName string, enableQuota bool, parameters map[string]string) error {

	// Get the space configuration without the space defaults loaded
	loadSpaceDefaults := false
	spaceConfig, err := m.GetASpaceConfig(orgName, spaceName, loadSpaceDefaults)
	if err != nil {
		return err
	}

	// Check to ensure that our target space config was found.
	if err != nil {
		return fmt.Errorf("Error retrieving space %s in Org %s because: %s", spaceName, orgName, err)
	} else if spaceConfig == nil {
		return fmt.Errorf("The space %s was not found in Org %s", spaceName, orgName)
	}

	// Set the Enable Space Quota Field
	spaceConfig.EnableSpaceQuota = enableQuota

	// For each parameter key (field), find it in the SpaceConfig struct
	// If the field is is found, update it with the parameter value.
	for field, value := range parameters {
		err := m.setFieldIn(spaceConfig, field, value)
		// We should always just error out if any of the parameters fail to set
		if err != nil {
			return err
		}
	}

	// Dump the file back out
	return m.UtilsMgr.WriteFile((*spaceConfig).GetSpaceConfigFilenameAndPath(m.ConfigDir, orgName, spaceName), *spaceConfig)
}

// UserMgmt specifies users and groups that can be associated to a particular org or space.
type UserMgmt struct {
	LDAPUsers  []string `yaml:"ldap_users"`
	Users      []string `yaml:"users"`
	SamlUsers  []string `yaml:"saml_users"`
	LDAPGroup  string   `yaml:"ldap_group"`
	LDAPGroups []string `yaml:"ldap_groups"`
}

func (u *UserMgmt) groups(groupName string) []string {
	groupMap := make(map[string]string)
	for _, group := range u.LDAPGroups {
		groupMap[group] = group
	}
	if u.LDAPGroup != "" {
		groupMap[u.LDAPGroup] = u.LDAPGroup
	}
	if groupName != "" {
		groupMap[groupName] = groupName
	}

	result := make([]string, 0, len(groupMap))
	for k := range groupMap {
		result = append(result, k)
	}
	return result
}
