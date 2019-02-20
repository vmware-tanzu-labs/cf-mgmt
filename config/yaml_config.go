package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/xchapter7x/lo"
)

const unlimited = "unlimited"

// yamlManager is the default implementation of Manager.
// It is backed by a directory of YAML files.
type yamlManager struct {
	ConfigDir string
}

// Orgs reads the config for all orgs.
func (m *yamlManager) Orgs() (*Orgs, error) {
	configFile := filepath.Join(m.ConfigDir, "orgs.yml")
	lo.G.Debug("Processing org file", configFile)
	input := &Orgs{}
	if err := LoadFile(configFile, &input); err != nil {
		return nil, err
	}
	return input, nil
}

func (m *yamlManager) GetDefaultASGConfigs() ([]ASGConfig, error) {
	filePath := path.Join(m.ConfigDir, "default_asgs")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		lo.G.Infof("No default asgs found.  Create directory default_asgs and add asg defintion(s)")
		return nil, nil
	}
	files, err := FindFiles(filePath, ".json")
	if err != nil {
		return nil, err
	}
	var result []ASGConfig
	for _, securityGroupFile := range files {
		lo.G.Debug("Loading security group contents", securityGroupFile)
		bytes, err := ioutil.ReadFile(securityGroupFile)
		if err != nil {
			return nil, err
		}
		asgConfig := ASGConfig{}
		lo.G.Debug("setting security group contents", string(bytes))
		asgConfig.Rules = string(bytes)
		asgConfig.Name = strings.Replace(filepath.Base(securityGroupFile), ".json", "", 1)
		result = append(result, asgConfig)
	}
	return result, nil
}

// GetASGConfigs reads all ASGs from the cf-mgmt configuration.
func (m *yamlManager) GetASGConfigs() ([]ASGConfig, error) {
	filePath := path.Join(m.ConfigDir, "asgs")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		lo.G.Infof("No asgs found.  Create directory asgs and add asg defintion(s)")
		return nil, nil
	}
	files, err := FindFiles(filePath, ".json")
	if err != nil {
		return nil, err
	}
	var result []ASGConfig
	for _, securityGroupFile := range files {
		lo.G.Debug("Loading security group contents", securityGroupFile)
		bytes, err := ioutil.ReadFile(securityGroupFile)
		if err != nil {
			return nil, err
		}
		asgConfig := ASGConfig{}
		lo.G.Debug("setting security group contents", string(bytes))
		asgConfig.Rules = string(bytes)
		asgConfig.Name = strings.Replace(filepath.Base(securityGroupFile), ".json", "", 1)
		result = append(result, asgConfig)
	}
	return result, nil
}

// GetIsolationSegmentConfig reads isolation segment config
func (m *yamlManager) GetGlobalConfig() (*GlobalConfig, error) {
	globalConfig := &GlobalConfig{}
	LoadFile(path.Join(m.ConfigDir, "cf-mgmt.yml"), globalConfig)
	return globalConfig, nil
}

// GetOrgConfigs reads all orgs from the cf-mgmt configuration.
func (m *yamlManager) GetOrgConfigs() ([]OrgConfig, error) {
	files, err := FindFiles(m.ConfigDir, "orgConfig.yml")
	if err != nil {
		return nil, err
	}
	result := make([]OrgConfig, len(files))
	for i, f := range files {
		result[i].AppTaskLimit = unlimited
		result[i].AppInstanceLimit = unlimited
		result[i].TotalReservedRoutePorts = "0"
		result[i].TotalPrivateDomains = unlimited
		result[i].TotalServiceKeys = unlimited
		result[i].InstanceMemoryLimit = unlimited

		if err = LoadFile(f, &result[i]); err != nil {
			lo.G.Error(err)
			return nil, err
		}
	}
	return result, nil
}

func (m *yamlManager) SaveOrgSpaces(spaces *Spaces) error {
	return WriteFile(filepath.Join(m.ConfigDir, spaces.Org, "spaces.yml"), spaces)
}

func (m *yamlManager) Spaces() ([]Spaces, error) {
	files, err := FindFiles(m.ConfigDir, "spaces.yml")
	if err != nil {
		return nil, err
	}

	spaceList := make([]Spaces, len(files))
	for i, f := range files {
		lo.G.Debug("Processing space file", f)

		if err = LoadFile(f, &spaceList[i]); err != nil {
			lo.G.Errorf("reading config for space %s: %v", f, err)
			return nil, err
		}
	}
	return spaceList, nil
}
func (m *yamlManager) OrgSpaces(orgName string) (*Spaces, error) {
	spaceList, err := m.Spaces()
	if err != nil {
		return nil, err
	}
	for _, space := range spaceList {
		if space.Org == orgName {
			return &space, nil
		}
	}
	return nil, fmt.Errorf("No spaces found for org [%s]", orgName)
}

func (m *yamlManager) GetSpaceConfigs() ([]SpaceConfig, error) {

	spaceDefaults := SpaceConfig{}
	LoadFile(filepath.Join(m.ConfigDir, "spaceDefaults.yml"), &spaceDefaults)

	files, err := FindFiles(m.ConfigDir, "spaceConfig.yml")
	if err != nil {
		return nil, err
	}
	result := make([]SpaceConfig, len(files))
	for i, f := range files {
		result[i].AppInstanceLimit = unlimited
		result[i].AppTaskLimit = unlimited
		result[i].TotalReservedRoutePorts = "0"
		result[i].TotalServiceKeys = unlimited
		result[i].InstanceMemoryLimit = unlimited

		if err = LoadFile(f, &result[i]); err != nil {
			return nil, err
		}

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

		if result[i].EnableSecurityGroup {
			securityGroupFile := strings.Replace(f, "spaceConfig.yml", "security-group.json", -1)
			lo.G.Debug("Loading security group contents", securityGroupFile)
			bytes, err := ioutil.ReadFile(securityGroupFile)
			if err != nil {
				return nil, err
			}
			lo.G.Debug("setting security group contents", string(bytes))
			result[i].SecurityGroupContents = string(bytes)
		}
	}
	return result, nil
}

func (m *yamlManager) GetOrgConfig(orgName string) (*OrgConfig, error) {
	configs, err := m.GetOrgConfigs()
	if err != nil {
		return nil, err
	}
	for _, config := range configs {
		if config.Org == orgName {
			return &config, nil
		}
	}
	return nil, fmt.Errorf("Org [%s] not found in config", orgName)
}

func (m *yamlManager) SaveOrgConfig(orgConfig *OrgConfig) error {
	directory := fmt.Sprintf("%s/%s", m.ConfigDir, orgConfig.Org)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return err
		}
	}

	return WriteFile(filepath.Join(m.ConfigDir, orgConfig.Org, "orgConfig.yml"), orgConfig)
}

func (m *yamlManager) RenameOrgConfig(orgConfig *OrgConfig) error {
	newDirectory := fmt.Sprintf("%s/%s", m.ConfigDir, orgConfig.Org)
	originalDirectory := fmt.Sprintf("%s/%s", m.ConfigDir, orgConfig.OriginalOrg)

	err := RenameDirectory(originalDirectory, newDirectory)
	if err != nil {
		return err
	}
	return m.SaveOrgConfig(orgConfig)
}

func (m *yamlManager) RenameSpaceConfig(spaceConfig *SpaceConfig) error {
	newDirectory := path.Join(m.ConfigDir, spaceConfig.Org, spaceConfig.Space)
	originalDirectory := path.Join(m.ConfigDir, spaceConfig.Org, spaceConfig.OriginalSpace)

	err := RenameDirectory(originalDirectory, newDirectory)
	if err != nil {
		return err
	}
	return m.SaveSpaceConfig(spaceConfig)
}

func (m *yamlManager) GetSpaceConfig(orgName, spaceName string) (*SpaceConfig, error) {
	targetPath := path.Join(m.ConfigDir, orgName, spaceName)
	files, err := FindFiles(targetPath, "spaceConfig.yml")
	if err != nil {
		return nil, fmt.Errorf("Space [%s] not found in org [%s] config", spaceName, orgName)
	}
	if len(files) != 1 {
		return nil, fmt.Errorf("Space [%s] not found in org [%s] config", spaceName, orgName)
	}

	result := &SpaceConfig{}
	if err = LoadFile(files[0], &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (m *yamlManager) SaveSpaceConfig(spaceConfig *SpaceConfig) error {
	if err := os.MkdirAll(fmt.Sprintf("%s/%s/%s", m.ConfigDir, spaceConfig.Org, spaceConfig.Space), 0755); err != nil {
		return err
	}
	return WriteFile(fmt.Sprintf("%s/%s/%s/spaceConfig.yml", m.ConfigDir, spaceConfig.Org, spaceConfig.Space), spaceConfig)
}

func (m *yamlManager) DeleteOrgConfig(orgName string) error {
	orgs, err := m.Orgs()
	if err != nil {
		return err
	}
	if orgs.Contains(orgName) {
		var orgList []string
		for _, org := range orgs.Orgs {
			if org != orgName {
				orgList = append(orgList, org)
			}
		}
		orgs.Orgs = orgList
		if err := m.SaveOrgs(orgs); err != nil {
			return err
		}
		os.RemoveAll(path.Join(m.ConfigDir, orgName))
	}
	return nil
}

func (m *yamlManager) DeleteSpaceConfig(orgName, spaceName string) error {
	spaces, err := m.OrgSpaces(orgName)
	if err != nil {
		return err
	}
	if spaces.Contains(spaceName) {
		var spaceList []string
		for _, space := range spaces.Spaces {
			if space != spaceName {
				spaceList = append(spaceList, space)
			}
		}
		spaces.Spaces = spaceList
		if err := m.SaveOrgSpaces(spaces); err != nil {
			return err
		}
		os.RemoveAll(path.Join(m.ConfigDir, orgName, spaceName))
	}
	return nil
}

// GetSpaceDefaults returns the default space configuration, if one was provided.
// If no space defaults were configured, a nil config and a nil error are returned.
func (m *yamlManager) GetSpaceDefaults() (*SpaceConfig, error) {
	fp := filepath.Join(m.ConfigDir, "spaceDefaults.yml")

	if !FileOrDirectoryExists(fp) {
		return nil, nil
	}
	result := SpaceConfig{}
	err := LoadFile(fp, &result)
	return &result, err
}

func (m *yamlManager) SaveOrgs(orgs *Orgs) error {
	if err := WriteFile(fmt.Sprintf("%s/orgs.yml", m.ConfigDir), orgs); err != nil {
		return err
	}
	return nil
}

// AddOrgToConfig adds an organization to the cf-mgmt configuration.
func (m *yamlManager) AddOrgToConfig(orgConfig *OrgConfig, spaces *Spaces) error {
	orgList, err := m.Orgs()
	if err != nil {
		return err
	}
	orgName := orgConfig.Org
	if orgName == "" {
		return errors.New("cannot have an empty org name")
	}

	if orgList.Contains(orgName) {
		return fmt.Errorf("org [%s] already added to config -> %v", orgName, orgList.Orgs)
	}
	lo.G.Infof("Adding org: %s ", orgName)
	orgList.Orgs = append(orgList.Orgs, orgName)
	if err = m.SaveOrgs(orgList); err != nil {
		return err
	}
	m.SaveOrgConfig(orgConfig)
	return m.SaveOrgSpaces(spaces)
}

// AddSpaceToConfig adds a space to the cf-mgmt configuration, so long as a
// space with the specified name doesn't already exist.
func (m *yamlManager) AddSpaceToConfig(spaceConfig *SpaceConfig) error {
	orgName := spaceConfig.Org
	spaceFileName := filepath.Join(m.ConfigDir, orgName, "spaces.yml")
	spaceList := &Spaces{}
	spaceName := spaceConfig.Space

	if err := LoadFile(spaceFileName, spaceList); err != nil {
		return err
	}
	if spaceList.Contains(spaceName) {
		return fmt.Errorf("space [%s] already added to config -> [%v]", spaceName, spaceList.Spaces)
	}
	lo.G.Infof("Adding space: %s ", spaceName)
	spaceList.Spaces = append(spaceList.Spaces, spaceName)
	if err := WriteFile(spaceFileName, spaceList); err != nil {
		return err
	}

	if err := m.SaveSpaceConfig(spaceConfig); err != nil {
		return err
	}
	if err := WriteFileBytes(fmt.Sprintf("%s/%s/%s/security-group.json", m.ConfigDir, orgName, spaceName), []byte("[]")); err != nil {
		return err
	}
	return nil
}

//AddSecurityGroupToSpace - adds security group json to org/space location
func (m *yamlManager) AddSecurityGroupToSpace(orgName, spaceName string, securityGroupDefinition []byte) error {
	return WriteFileBytes(fmt.Sprintf("%s/%s/%s/security-group.json", m.ConfigDir, orgName, spaceName), securityGroupDefinition)
}

//AddSecurityGroupToSpace - adds security group json to org/space location
func (m *yamlManager) AddSecurityGroup(securityGroupName string, securityGroupDefinition []byte) error {
	lo.G.Infof("Writing out bytes for security group %s", securityGroupName)
	return WriteFileBytes(fmt.Sprintf("%s/asgs/%s.json", m.ConfigDir, securityGroupName), securityGroupDefinition)
}

//AddDefaultSecurityGroup - adds security group json to org/space location
func (m *yamlManager) AddDefaultSecurityGroup(securityGroupName string, securityGroupDefinition []byte) error {
	lo.G.Infof("Writing out bytes for security group %s", securityGroupName)
	return WriteFileBytes(fmt.Sprintf("%s/default_asgs/%s.json", m.ConfigDir, securityGroupName), securityGroupDefinition)
}

func (m *yamlManager) AddOrgQuota(orgQuota OrgQuota) error {
	lo.G.Infof("Writing out orgQuota %s", orgQuota.Name)
	return WriteFile(fmt.Sprintf("%s/org_quotas/%s.yml", m.ConfigDir, orgQuota.Name), orgQuota)
}

func (m *yamlManager) AddSpaceQuota(spaceQuota SpaceQuota) error {
	quotasDir := path.Join(m.ConfigDir, spaceQuota.Org, "space_quotas")
	if err := os.MkdirAll(quotasDir, 0755); err != nil {
		lo.G.Errorf("Error creating config directory %s. Error : %s", quotasDir, err)
		return fmt.Errorf("cannot create directory %s: %v", quotasDir, err)
	}
	lo.G.Infof("Writing out spaceQuota %s for org %s", spaceQuota.Name, spaceQuota.Org)
	return WriteFile(fmt.Sprintf("%s/%s/space_quotas/%s.yml", m.ConfigDir, spaceQuota.Org, spaceQuota.Name), spaceQuota)
}

// CreateConfigIfNotExists initializes a new configuration directory.
// If the specified configuration directory already exists, it is left unmodified.
func (m *yamlManager) CreateConfigIfNotExists(uaaOrigin string) error {
	if FileOrDirectoryExists(m.ConfigDir) {
		lo.G.Infof("Config directory %s already exists, skipping creation", m.ConfigDir)
		return nil
	}
	if err := os.MkdirAll(m.ConfigDir, 0755); err != nil {
		lo.G.Errorf("Error creating config directory %s. Error : %s", m.ConfigDir, err)
		return fmt.Errorf("cannot create directory %s: %v", m.ConfigDir, err)
	}
	lo.G.Infof("Config directory %s created", m.ConfigDir)

	asgDir := path.Join(m.ConfigDir, "asgs")
	if err := os.MkdirAll(asgDir, 0755); err != nil {
		lo.G.Errorf("Error creating config directory %s. Error : %s", asgDir, err)
		return fmt.Errorf("cannot create directory %s: %v", asgDir, err)
	}
	lo.G.Infof("ASG directory %s created", asgDir)

	asgDir = path.Join(m.ConfigDir, "default_asgs")
	if err := os.MkdirAll(asgDir, 0755); err != nil {
		lo.G.Errorf("Error creating config directory %s. Error : %s", asgDir, err)
		return fmt.Errorf("cannot create directory %s: %v", asgDir, err)
	}
	lo.G.Infof("ASG directory %s created", asgDir)

	orgQuotasDir := path.Join(m.ConfigDir, "org_quotas")
	if err := os.MkdirAll(orgQuotasDir, 0755); err != nil {
		lo.G.Errorf("Error creating config directory %s. Error : %s", orgQuotasDir, err)
		return fmt.Errorf("cannot create directory %s: %v", orgQuotasDir, err)
	}
	lo.G.Infof("OrgQuotas directory %s created", orgQuotasDir)

	if err := m.SaveGlobalConfig(&GlobalConfig{}); err != nil {
		return err
	}
	if err := WriteFile(fmt.Sprintf("%s/ldap.yml", m.ConfigDir), &LdapConfig{TLS: false, Origin: uaaOrigin}); err != nil {
		return err
	}

	if err := WriteFile(fmt.Sprintf("%s/orgs.yml", m.ConfigDir), &Orgs{
		EnableDeleteOrgs: true,
		ProtectedOrgs:    DefaultProtectedOrgs,
	}); err != nil {
		return err
	}
	if err := WriteFile(fmt.Sprintf("%s/spaceDefaults.yml", m.ConfigDir), struct {
		Developer UserMgmt `yaml:"space-developer"`
		Manager   UserMgmt `yaml:"space-manager"`
		Auditor   UserMgmt `yaml:"space-auditor"`
	}{}); err != nil {
		return err
	}
	return nil
}

func (m *yamlManager) SaveGlobalConfig(globalConfig *GlobalConfig) error {
	return WriteFile(fmt.Sprintf("%s/cf-mgmt.yml", m.ConfigDir), globalConfig)
}

// DeleteConfigIfExists deletes config directory if it exists.
func (m *yamlManager) DeleteConfigIfExists() error {
	if !FileOrDirectoryExists(m.ConfigDir) {
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

func (m *yamlManager) LdapConfig(ldapBindPassword string) (*LdapConfig, error) {
	config := &LdapConfig{}
	err := LoadFile(path.Join(m.ConfigDir, "ldap.yml"), config)
	if err != nil {
		return nil, err
	}
	if ldapBindPassword != "" {
		config.BindPassword = ldapBindPassword
	} else {
		lo.G.Warning("Ldap bind password should be removed from ldap.yml as this will be deprecated in a future release.  Use --ldap-password flag instead.")
	}
	if config.Origin == "" {
		config.Origin = "ldap"
	}
	return config, nil
}

func (m *yamlManager) GetOrgQuotas() ([]OrgQuota, error) {
	filePath := path.Join(m.ConfigDir, "org_quotas")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		lo.G.Infof("No org quotas found.  Create directory org_quotas and add org quota defintion(s)")
		return nil, nil
	}
	files, err := FindFiles(filePath, ".yml")
	if err != nil {
		return nil, err
	}
	var result []OrgQuota
	for _, orgQuotaFile := range files {
		orgQuota := &OrgQuota{}
		err = LoadFile(orgQuotaFile, orgQuota)
		if err != nil {
			return nil, err
		}
		orgQuota.Name = strings.Replace(filepath.Base(orgQuotaFile), ".yml", "", 1)
		result = append(result, *orgQuota)
	}
	return result, nil
}

func (m *yamlManager) GetOrgQuota(name string) (*OrgQuota, error) {
	orgQuotas, err := m.GetOrgQuotas()
	if err != nil {
		return nil, err
	}
	for _, orgQuota := range orgQuotas {
		if strings.EqualFold(orgQuota.Name, name) {
			return &orgQuota, nil
		}
	}
	return nil, nil
}

func (m *yamlManager) SaveOrgQuota(orgQuota *OrgQuota) error {
	orgQuotaPath := path.Join(m.ConfigDir, "org_quotas")
	if err := os.MkdirAll(orgQuotaPath, 0755); err != nil {
		return fmt.Errorf("cannot create directory %s: %v", orgQuotaPath, err)
	}
	fmt.Println(fmt.Sprintf("Saving Named Org Quote %s", orgQuota.Name))
	return WriteFile(fmt.Sprintf("%s/%s.yml", orgQuotaPath, orgQuota.Name), orgQuota)
}

func (m *yamlManager) GetSpaceQuotas(org string) ([]SpaceQuota, error) {
	filePath := path.Join(m.ConfigDir, org, "space_quotas")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		lo.G.Infof("No space quotas found. Create directory space_quotas for org %s and add space quota defintion(s)", org)
		return nil, nil
	}
	files, err := FindFiles(filePath, ".yml")
	if err != nil {
		return nil, err
	}
	var result []SpaceQuota
	for _, spaceQuotaFile := range files {
		spaceQuota := &SpaceQuota{}
		err = LoadFile(spaceQuotaFile, spaceQuota)
		if err != nil {
			return nil, err
		}
		spaceQuota.Name = strings.Replace(filepath.Base(spaceQuotaFile), ".yml", "", 1)
		spaceQuota.Org = org
		result = append(result, *spaceQuota)
	}
	return result, nil
}

func (m *yamlManager) GetSpaceQuota(name, org string) (*SpaceQuota, error) {
	spaceQuotas, err := m.GetSpaceQuotas(org)
	if err != nil {
		return nil, err
	}
	for _, spaceQuota := range spaceQuotas {
		if strings.EqualFold(spaceQuota.Name, name) {
			return &spaceQuota, nil
		}
	}
	return nil, nil
}

func (m *yamlManager) SaveSpaceQuota(spaceQuota *SpaceQuota) error {
	spaceQuotaPath := path.Join(m.ConfigDir, spaceQuota.Org, "space_quotas")
	if err := os.MkdirAll(spaceQuotaPath, 0755); err != nil {
		return fmt.Errorf("cannot create directory %s: %v", spaceQuotaPath, err)
	}
	targetFile := fmt.Sprintf("%s/%s.yml", spaceQuotaPath, spaceQuota.Name)
	fmt.Println(fmt.Sprintf("Saving Named Space Quote %s for org %s", spaceQuota.Name, spaceQuota.Org))
	return WriteFile(targetFile, spaceQuota)
}
