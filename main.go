package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/export"
	"github.com/pivotalservices/cf-mgmt/generated"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/xchapter7x/lo"
)

var (
	//VERSION -
	VERSION string
)

//ErrorHandler -
type ErrorHandler struct {
	ExitCode int
	Error    error
}

type flagBucket struct {
	Desc        string
	EnvVar      string
	StringSlice bool
}

//CFMgmt -
type CFMgmt struct {
	UAAManager      uaa.Manager
	OrgManager      organization.Manager
	SpaceManager    space.Manager
	ConfigManager   config.Manager
	ConfigDirectory string
	PeekDeletion    bool
	LdapBindPwd     string
	UaacToken       string
	SystemDomain    string
	UAACManager     uaac.Manager
	CloudController cloudcontroller.Manager
}

//InitializeManager -
func InitializeManager(c *cli.Context) (*CFMgmt, error) {
	var err error
	configDir := getConfigDir(c)
	sysDomain := c.String(getFlag(systemDomain))
	user := c.String(getFlag(userID))
	pwd := c.String(getFlag(password))
	secret := c.String(getFlag(clientSecret))
	ldapPwd := c.String(getFlag(ldapPassword))
	peek := c.Bool("peek")

	if sysDomain == "" ||
		user == "" ||
		secret == "" {
		err = fmt.Errorf("must set system-domain, user-id, client-secret properties")
		return nil, err
	}

	var cfToken, uaacToken string
	cfMgmt := &CFMgmt{}
	cfMgmt.LdapBindPwd = ldapPwd
	cfMgmt.PeekDeletion = peek
	cfMgmt.ConfigDirectory = configDir
	cfMgmt.SystemDomain = sysDomain
	cfMgmt.UAAManager = uaa.NewDefaultUAAManager(sysDomain, user)

	if uaacToken, err = cfMgmt.UAAManager.GetUAACToken(secret); err != nil {
		return nil, err
	}
	cfMgmt.UaacToken = uaacToken
	cfMgmt.UAACManager = uaac.NewManager(sysDomain, uaacToken)

	if pwd != "" {
		lo.G.Warning("Password parameter is deprecated, create uaa client and client-secret instead")
		if cfToken, err = cfMgmt.UAAManager.GetCFToken(pwd); err != nil {
			return nil, err
		}
		cfMgmt.CloudController = cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain), cfToken)
	} else {
		cfToken = uaacToken
		cfMgmt.CloudController = cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", cfMgmt.SystemDomain), uaacToken)
	}
	cfMgmt.OrgManager = organization.NewManager(sysDomain, cfToken, uaacToken)
	cfMgmt.SpaceManager = space.NewManager(sysDomain, cfToken, uaacToken)
	cfMgmt.ConfigManager = config.NewManager(configDir)

	return cfMgmt, nil
}

const (
	systemDomain     string = "SYSTEM_DOMAIN"
	userID           string = "USER_ID"
	password         string = "PASSWORD"
	clientSecret     string = "CLIENT_SECRET"
	configDirectory  string = "CONFIG_DIR"
	orgName          string = "ORG"
	spaceName        string = "SPACE"
	ldapPassword     string = "LDAP_PASSWORD"
	orgBillingMgrGrp string = "ORG_BILLING_MGR_GRP"
	orgMgrGrp        string = "ORG_MGR_GRP"
	orgAuditorGrp    string = "ORG_AUDITOR_GRP"
	spaceDevGrp      string = "SPACE_DEV_GRP"
	spaceMgrGrp      string = "SPACE_MGR_GRP"
	spaceAuditorGrp  string = "SPACE_AUDITOR_GRP"
)

func main() {
	eh := new(ErrorHandler)
	eh.ExitCode = 0
	app := NewApp(eh)
	if err := app.Run(os.Args); err != nil {
		eh.ExitCode = 1
		eh.Error = err
		lo.G.Error(eh.Error)
	}
	os.Exit(eh.ExitCode)
}

// NewApp creates a new cli app
func NewApp(eh *ErrorHandler) *cli.App {
	//cli.AppHelpTemplate = CfopsHelpTemplate
	app := cli.NewApp()
	app.Version = VERSION
	app.Name = "cf-mgmt"
	app.Usage = "cf-mgmt"
	app.Commands = []cli.Command{
		{
			Name:  "version",
			Usage: "shows the application version currently in use",
			Action: func(c *cli.Context) (err error) {
				cli.ShowVersion(c)
				return
			},
		},
		CreateInitCommand(eh),
		CreateAddOrgCommand(eh),
		CreateAddSpaceCommand(eh),
		CreateExportConfigCommand(eh),
		CreateGeneratePipelineCommand(runGeneratePipeline, eh),
		CreateCommand("create-orgs", runCreateOrgs, defaultFlags(), eh),
		CreateCommand("create-org-private-domains", runCreateOrgPrivateDomains, defaultFlags(), eh),
		CreateCommand("delete-orgs", runDeleteOrgs, defaultFlagsWithDelete(), eh),
		CreateCommand("update-org-quotas", runCreateOrgQuotas, defaultFlags(), eh),
		CreateCommand("update-org-users", runUpdateOrgUsers, defaultFlagsWithLdap(), eh),
		CreateCommand("create-spaces", runCreateSpaces, defaultFlagsWithLdap(), eh),
		CreateCommand("delete-spaces", runDeleteSpaces, defaultFlagsWithDelete(), eh),
		CreateCommand("update-spaces", runUpdateSpaces, defaultFlags(), eh),
		CreateCommand("update-space-quotas", runCreateSpaceQuotas, defaultFlags(), eh),
		CreateCommand("update-space-users", runUpdateSpaceUsers, defaultFlagsWithLdap(), eh),
		CreateCommand("update-space-security-groups", runCreateSpaceSecurityGroups, defaultFlags(), eh),
	}

	return app
}

// CreateExportConfigCommand - Creates CLI command for export config
func CreateExportConfigCommand(eh *ErrorHandler) (command cli.Command) {
	flags := defaultFlags()
	flag := cli.StringSliceFlag{
		Name:  "excluded-org",
		Usage: "Org to be excluded from export. Repeat the flag to specify multiple orgs",
	}
	flags = append(flags, flag)
	flag = cli.StringSliceFlag{
		Name:  "excluded-space",
		Usage: "Space to be excluded from export. Repeat the flag to specify multiple spaces",
	}
	flags = append(flags, flag)
	command = cli.Command{
		Name:        "export-config",
		Usage:       "Exports org, space and user configuration from an existing CF instance. Try export-config --help for more options",
		Description: "Exports org and space configurations from an existing Cloud Foundry instance. [Warning: This operation will delete existing config folder]",
		Action:      runExportConfig,
		Flags:       flags,
	}
	return
}

//CreateInitCommand -
func CreateInitCommand(eh *ErrorHandler) (command cli.Command) {
	flagList := map[string]flagBucket{
		configDirectory: {
			Desc:   "Name of the config directory. Default config directory is `config`",
			EnvVar: configDirectory,
		},
	}

	command = cli.Command{
		Name:        "init-config",
		Usage:       "Initializes folder structure for configuration",
		Description: "Initializes folder structure for configuration",
		Action:      runInit,
		Flags:       buildFlags(flagList),
	}
	return
}

func runInit(c *cli.Context) (err error) {
	configDir := getConfigDir(c)
	configManager := config.NewManager(configDir)
	err = configManager.CreateConfigIfNotExists("ldap")
	return err
}

//CreateAddOrgCommand -
func CreateAddOrgCommand(eh *ErrorHandler) (command cli.Command) {
	flagList := map[string]flagBucket{
		configDirectory: {
			Desc:   "Config directory name.  Default is config",
			EnvVar: configDirectory,
		},
		orgName: {
			Desc:   "Org name to add",
			EnvVar: orgName,
		},
		orgBillingMgrGrp: {
			Desc:   "LDAP group for Org Billing Manager",
			EnvVar: orgBillingMgrGrp,
		},
		orgMgrGrp: {
			Desc:   "LDAP group for Org Manager",
			EnvVar: orgMgrGrp,
		},
		orgAuditorGrp: {
			Desc:   "LDAP group for Org Auditor",
			EnvVar: orgAuditorGrp,
		},
	}

	command = cli.Command{
		Name:        "add-org-to-config",
		Usage:       "Adds specified org to configuration",
		Description: "Adds specified org to configuration",
		Action:      runAddOrg,
		Flags:       buildFlags(flagList),
	}
	return
}

func runAddOrg(c *cli.Context) error {
	inputOrg := c.String(getFlag(orgName))
	orgConfig := &config.OrgConfig{OrgName: inputOrg,
		OrgBillingMgrLDAPGrp: c.String(getFlag(orgBillingMgrGrp)),
		OrgMgrLDAPGrp:        c.String(getFlag(orgMgrGrp)),
		OrgAuditorLDAPGrp:    c.String(getFlag(orgAuditorGrp)),
	}
	return config.NewManager(getConfigDir(c)).AddOrgToConfig(orgConfig)
}

//CreateAddSpaceCommand -
func CreateAddSpaceCommand(eh *ErrorHandler) (command cli.Command) {
	flagList := map[string]flagBucket{
		configDirectory: {
			Desc:   "config dir.  Default is config",
			EnvVar: configDirectory,
		},
		orgName: {
			Desc:   "org name of space",
			EnvVar: orgName,
		},
		spaceName: {
			Desc:   "space name to add",
			EnvVar: spaceName,
		},
		spaceDevGrp: {
			Desc:   "LDAP group for Space Developer",
			EnvVar: spaceDevGrp,
		},
		spaceMgrGrp: {
			Desc:   "LDAP group for Space Manager",
			EnvVar: spaceMgrGrp,
		},
		spaceAuditorGrp: {
			Desc:   "LDAP group for Space Auditor",
			EnvVar: spaceAuditorGrp,
		},
	}

	command = cli.Command{
		Name:        "add-space-to-config",
		Usage:       "adds specified space to configuration for org",
		Description: "adds specified space to configuration for org",
		Action:      runAddSpace,
		Flags:       buildFlags(flagList),
	}
	return
}

func runAddSpace(c *cli.Context) (err error) {

	inputOrg := c.String(getFlag(orgName))
	inputSpace := c.String(getFlag(spaceName))

	spaceConfig := &config.SpaceConfig{OrgName: inputOrg,
		SpaceName:           inputSpace,
		SpaceDevLDAPGrp:     c.String(getFlag(spaceDevGrp)),
		SpaceMgrLDAPGrp:     c.String(getFlag(spaceMgrGrp)),
		SpaceAuditorLDAPGrp: c.String(getFlag(spaceAuditorGrp)),
	}

	configDr := getConfigDir(c)
	if inputOrg == "" || inputSpace == "" {
		err = fmt.Errorf("Must provide org name and space name")
	} else {
		err = config.NewManager(configDr).AddSpaceToConfig(spaceConfig)
	}
	return
}

//CreateGeneratePipelineCommand -
func CreateGeneratePipelineCommand(action func(c *cli.Context) (err error), eh *ErrorHandler) (command cli.Command) {
	command = cli.Command{
		Name:        "generate-concourse-pipeline",
		Usage:       "generates a concourse pipline based on convention of org/space metadata",
		Description: "generate-concourse-pipeline",
		Action:      action,
	}
	return
}

func runGeneratePipeline(c *cli.Context) (err error) {
	const varsFileName = "vars.yml"
	const pipelineFileName = "pipeline.yml"
	const cfMgmtYml = "cf-mgmt.yml"
	const cfMgmtSh = "cf-mgmt.sh"
	var targetFile string
	fmt.Println("Generating pipeline....")
	if err = createFile(pipelineFileName, pipelineFileName); err != nil {
		lo.G.Error("Error creating pipeline.yml", err)
		return
	}
	if err = createFile(varsFileName, varsFileName); err != nil {
		lo.G.Error("Error creating vars.yml", err)
		return
	}
	if err = os.MkdirAll("ci/tasks", 0755); err == nil {
		targetFile = fmt.Sprintf("./ci/tasks/%s", cfMgmtYml)
		lo.G.Debug("Creating", targetFile)
		if err = createFile(cfMgmtYml, targetFile); err != nil {
			lo.G.Error("Error creating cf-mgmt.yml", err)
			return
		}
		targetFile = fmt.Sprintf("./ci/tasks/%s", cfMgmtSh)
		lo.G.Debug("Creating", targetFile)
		if err = createFile(cfMgmtSh, targetFile); err != nil {
			lo.G.Error("Error creating cf-mgmt.sh", err)
			return
		}
	}
	fmt.Println("1) Update vars.yml with the appropriate values")
	fmt.Println("2) Using following command to set your pipeline in concourse after you have checked all files in to GIT")
	fmt.Println("fly -t lite set-pipeline -p cf-mgmt -c pipeline.yml --load-vars-from=vars.yml")
	return
}

func createFile(assetName, fileName string) error {
	bytes, err := generated.Asset(fmt.Sprintf("files/%s", assetName))
	if err != nil {
		return err
	}
	perm := os.FileMode(0666)
	if strings.HasSuffix(fileName, ".sh") {
		perm = 0755
	}
	return ioutil.WriteFile(fileName, bytes, perm)
}

//CreateCommand -
func CreateCommand(commandName string, action func(c *cli.Context) (err error), flags []cli.Flag, eh *ErrorHandler) (command cli.Command) {
	command = cli.Command{
		Name:        commandName,
		Usage:       fmt.Sprintf("%s with what is defined in config", commandName),
		Description: commandName,
		Action:      action,
		Flags:       flags,
	}
	return
}

func runCreateOrgs(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.OrgManager.CreateOrgs(cfMgmt.ConfigDirectory)
	}
	return err
}

func runCreateOrgPrivateDomains(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.OrgManager.CreatePrivateDomains(cfMgmt.ConfigDirectory)
	}
	return err
}

func runDeleteOrgs(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.OrgManager.DeleteOrgs(cfMgmt.ConfigDirectory, cfMgmt.PeekDeletion)
	}
	return err
}

func runDeleteSpaces(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.DeleteSpaces(cfMgmt.ConfigDirectory, cfMgmt.PeekDeletion)
	}
	return err
}

func runCreateOrgQuotas(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.OrgManager.CreateQuotas(cfMgmt.ConfigDirectory)
	}
	return err
}

func runCreateSpaceQuotas(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.CreateQuotas(cfMgmt.ConfigDirectory)
	}
	return err
}

func runCreateSpaceSecurityGroups(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.CreateApplicationSecurityGroups(cfMgmt.ConfigDirectory)
	}
	return err
}

func runCreateSpaces(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.CreateSpaces(cfMgmt.ConfigDirectory, cfMgmt.LdapBindPwd)
	}
	return err
}

func runUpdateSpaces(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.UpdateSpaces(cfMgmt.ConfigDirectory)
	}
	return err
}

func runUpdateSpaceUsers(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.UpdateSpaceUsers(cfMgmt.ConfigDirectory, cfMgmt.LdapBindPwd)
	}
	return err
}

func runUpdateOrgUsers(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.OrgManager.UpdateOrgUsers(cfMgmt.ConfigDirectory, cfMgmt.LdapBindPwd)
	}
	return err
}

func defaultFlagsWithLdap() (flags []cli.Flag) {
	flags = defaultFlags()
	flag := cli.StringFlag{
		Name:   getFlag(ldapPassword),
		EnvVar: ldapPassword,
		Usage:  "Ldap password for binding",
	}
	flags = append(flags, flag)
	return
}

func defaultFlagsWithDelete() (flags []cli.Flag) {
	flags = defaultFlags()
	flag := cli.BoolFlag{
		Name:  "peek",
		Usage: "Preview entities to delete without deleting them.",
	}
	flags = append(flags, flag)
	return
}

func defaultFlags() (flags []cli.Flag) {
	var flagList = buildDefaultFlags()
	flags = buildFlags(flagList)
	return
}

func buildDefaultFlags() (flagList map[string]flagBucket) {
	flagList = map[string]flagBucket{
		systemDomain: {
			Desc:   "system domain",
			EnvVar: systemDomain,
		},
		userID: {
			Desc:   "user id that has privileges to create/update/delete users, orgs and spaces",
			EnvVar: userID,
		},
		password: {
			Desc:   "password for user account [optional if client secret is provided]",
			EnvVar: password,
		},
		clientSecret: {
			Desc:   "secret for user account that has sufficient privileges to create/update/delete users, orgs and spaces]",
			EnvVar: clientSecret,
		},
		configDirectory: {
			Desc:   "config dir.  Default is config",
			EnvVar: configDirectory,
		},
	}
	return
}
func buildFlags(flagList map[string]flagBucket) (flags []cli.Flag) {
	for _, v := range flagList {
		if v.StringSlice {
			flags = append(flags, cli.StringSliceFlag{
				Name:   getFlag(v.EnvVar),
				Usage:  v.Desc,
				EnvVar: v.EnvVar,
			})
		} else {
			flags = append(flags, cli.StringFlag{
				Name:   getFlag(v.EnvVar),
				Value:  "",
				Usage:  v.Desc,
				EnvVar: v.EnvVar,
			})
		}
	}
	return
}

func getFlag(input string) string {
	return strings.ToLower(strings.Replace(input, "_", "-", -1))
}

func getConfigDir(c *cli.Context) (cDir string) {
	cDir = c.String(getFlag(configDirectory))
	if cDir == "" {
		return "config"
	}
	return cDir
}

func runExportConfig(c *cli.Context) error {
	var cfMgmt *CFMgmt
	var err error
	cfMgmt, err = InitializeManager(c)
	if cfMgmt != nil {
		exportManager := export.NewExportManager(cfMgmt.ConfigDirectory, cfMgmt.UAACManager, cfMgmt.CloudController)
		excludedOrgs := make(map[string]string)
		excludedOrgs["system"] = "system"
		orgsExcludedByUser := c.StringSlice(getFlag("excluded-org"))
		for _, org := range orgsExcludedByUser {
			excludedOrgs[org] = org
		}
		excludedSpaces := make(map[string]string)
		spacesExcludedByUser := c.StringSlice(getFlag("excluded-space"))
		for _, space := range spacesExcludedByUser {
			excludedSpaces[space] = space
		}
		lo.G.Info("Orgs excluded from export by default: [system]")
		lo.G.Infof("Orgs excluded from export by user:  %v ", orgsExcludedByUser)
		lo.G.Infof("Spaces excluded from export by user:  %v ", spacesExcludedByUser)
		err = exportManager.ExportConfig(excludedOrgs, excludedSpaces)
		if err != nil {
			lo.G.Errorf("Export failed with error:  %s", err)
		}
	} else {
		lo.G.Errorf("Unable to initialize cf-mgmt. Error : %s", err)
	}
	return err
}
