package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/generated"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/pivotalservices/cf-mgmt/utils"
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
	UAAManager   uaa.Manager
	OrgManager   organization.Manager
	SpaceManager space.Manager
	ConfigDir    string
	LdapBindPwd  string
}

//InitializeManager -
func InitializeManager(c *cli.Context) (*CFMgmt, error) {
	var err error
	sysDomain := c.String(getFlag(systemDomain))
	user := c.String(getFlag(userID))
	pwd := c.String(getFlag(password))
	config := getConfigDir(c)
	secret := c.String(getFlag(clientSecret))
	ldapPwd := c.String(getFlag(ldapPassword))

	if sysDomain == "" ||
		user == "" ||
		pwd == "" ||
		secret == "" {
		err = fmt.Errorf("Must set system-domain, user-id, password, client-secret properties")
		return nil, err
	}

	var cfToken, uaacToken string
	cfMgmt := &CFMgmt{}
	cfMgmt.LdapBindPwd = ldapPwd
	cfMgmt.UAAManager = uaa.NewDefaultUAAManager(sysDomain, user)
	if cfToken, err = cfMgmt.UAAManager.GetCFToken(pwd); err != nil {
		return nil, err
	}
	if uaacToken, err = cfMgmt.UAAManager.GetUAACToken(secret); err != nil {
		return nil, err
	}
	cfMgmt.OrgManager = organization.NewManager(sysDomain, cfToken, uaacToken)
	cfMgmt.SpaceManager = space.NewManager(sysDomain, cfToken, uaacToken)
	cfMgmt.ConfigDir = config
	return cfMgmt, nil

}

const (
	systemDomain string = "SYSTEM_DOMAIN"
	userID       string = "USER_ID"
	password     string = "PASSWORD"
	clientSecret string = "CLIENT_SECRET"
	configDir    string = "CONFIG_DIR"
	orgName      string = "ORG"
	spaceName    string = "SPACE"
	ldapPassword string = "LDAP_PASSWORD"
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
		cli.Command{
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
		CreateGeneratePipelineCommand(runGeneratePipeline, eh),
		CreateCommand("create-orgs", runCreateOrgs, defaultFlags(), eh),
		CreateCommand("update-org-quotas", runCreateOrgQuotas, defaultFlags(), eh),
		CreateCommand("update-org-users", runUpdateOrgUsers, defaultFlagsWithLdap(), eh),
		CreateCommand("create-spaces", runCreateSpaces, defaultFlagsWithLdap(), eh),
		CreateCommand("update-spaces", runUpdateSpaces, defaultFlags(), eh),
		CreateCommand("update-space-quotas", runCreateSpaceQuotas, defaultFlags(), eh),
		CreateCommand("update-space-users", runUpdateSpaceUsers, defaultFlagsWithLdap(), eh),
		CreateCommand("update-space-security-groups", runCreateSpaceSecurityGroups, defaultFlags(), eh),
	}

	return app
}

//CreateInitCommand -
func CreateInitCommand(eh *ErrorHandler) (command cli.Command) {
	flagList := map[string]flagBucket{
		configDir: flagBucket{
			Desc:   "config dir.  Default is config",
			EnvVar: configDir,
		},
	}

	command = cli.Command{
		Name:        "init-config",
		Usage:       "initializes folder structure for configuration",
		Description: "initializes folder structure for configuration",
		Action:      runInit,
		Flags:       buildFlags(flagList),
	}
	return
}

func runInit(c *cli.Context) (err error) {
	config := getConfigDir(c)
	if err = os.MkdirAll(config, 0755); err == nil {
		utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/ldap.yml", config), &ldap.Config{TLS: false, Origin: "ldap"})
		utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/orgs.yml", config), &organization.InputOrgs{})
		utils.NewDefaultManager().WriteFile(fmt.Sprintf("%s/spaceDefaults.yml", config), &space.ConfigSpaceDefaults{})
	}
	return
}

//CreateAddOrgCommand -
func CreateAddOrgCommand(eh *ErrorHandler) (command cli.Command) {
	flagList := map[string]flagBucket{
		configDir: flagBucket{
			Desc:   "config dir.  Default is config",
			EnvVar: configDir,
		},
		orgName: flagBucket{
			Desc:   "org name to add",
			EnvVar: orgName,
		},
	}

	command = cli.Command{
		Name:        "add-org-to-config",
		Usage:       "adds specified org to configuration",
		Description: "adds specified org to configuration",
		Action:      runAddOrg,
		Flags:       buildFlags(flagList),
	}
	return
}

func runAddOrg(c *cli.Context) (err error) {
	inputOrg := c.String(getFlag(orgName))
	configDir := getConfigDir(c)
	if inputOrg == "" {
		err = fmt.Errorf("Must provide org name")
	} else {
		err = config.NewManager(configDir).AddOrgToConfig(inputOrg)
	}
	return
}

//CreateAddSpaceCommand -
func CreateAddSpaceCommand(eh *ErrorHandler) (command cli.Command) {
	flagList := map[string]flagBucket{
		configDir: flagBucket{
			Desc:   "config dir.  Default is config",
			EnvVar: configDir,
		},
		orgName: flagBucket{
			Desc:   "org name of space",
			EnvVar: orgName,
		},
		spaceName: flagBucket{
			Desc:   "space name to add",
			EnvVar: spaceName,
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
	configDir := getConfigDir(c)
	if inputOrg == "" || inputSpace == "" {
		err = fmt.Errorf("Must provide org name and space name")
	} else {
		err = config.NewManager(configDir).AddSpaceToConfig(inputOrg, inputSpace)
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

func createFile(assetName, fileName string) (err error) {
	var f *os.File
	var fileBytes []byte
	if fileBytes, err = generated.Asset(fmt.Sprintf("files/%s", assetName)); err == nil {
		if f, err = os.Create(fileName); err == nil {
			defer f.Close()
			_, err = f.Write(fileBytes)
		}
	}
	return
}

//CreateCommand -
func CreateCommand(commandName string, action func(c *cli.Context) (err error), flags []cli.Flag, eh *ErrorHandler) (command cli.Command) {
	desc := fmt.Sprintf(commandName)
	command = cli.Command{
		Name:        commandName,
		Usage:       fmt.Sprintf("%s with what is defined in config", commandName),
		Description: desc,
		Action:      action,
		Flags:       flags,
	}
	return
}

func runCreateOrgs(c *cli.Context) (err error) {
	var cfMgmt *CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.OrgManager.CreateOrgs(cfMgmt.ConfigDir)
	}
	return
}

func runCreateOrgQuotas(c *cli.Context) (err error) {
	var cfMgmt *CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.OrgManager.CreateQuotas(cfMgmt.ConfigDir)
	}
	return
}

func runCreateSpaceQuotas(c *cli.Context) (err error) {
	var cfMgmt *CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.CreateQuotas(cfMgmt.ConfigDir)
	}
	return
}

func runCreateSpaceSecurityGroups(c *cli.Context) (err error) {
	var cfMgmt *CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.CreateApplicationSecurityGroups(cfMgmt.ConfigDir)
	}
	return
}

func runCreateSpaces(c *cli.Context) (err error) {
	var cfMgmt *CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.CreateSpaces(cfMgmt.ConfigDir, cfMgmt.LdapBindPwd)
	}
	return
}

func runUpdateSpaces(c *cli.Context) (err error) {
	var cfMgmt *CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.UpdateSpaces(cfMgmt.ConfigDir)
	}
	return
}

func runUpdateSpaceUsers(c *cli.Context) (err error) {
	var cfMgmt *CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.UpdateSpaceUsers(cfMgmt.ConfigDir, cfMgmt.LdapBindPwd)
	}
	return
}

func runUpdateOrgUsers(c *cli.Context) (err error) {
	var cfMgmt *CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.OrgManager.UpdateOrgUsers(cfMgmt.ConfigDir, cfMgmt.LdapBindPwd)
	}
	return
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

func defaultFlags() (flags []cli.Flag) {
	var flagList = buildDefaultFlags()
	flags = buildFlags(flagList)
	return
}

func buildDefaultFlags() (flagList map[string]flagBucket) {
	flagList = map[string]flagBucket{
		systemDomain: flagBucket{
			Desc:   "system domain",
			EnvVar: systemDomain,
		},
		userID: flagBucket{
			Desc:   "user id that has admin priv",
			EnvVar: userID,
		},
		password: flagBucket{
			Desc:   "password for user account that has admin priv",
			EnvVar: password,
		},
		clientSecret: flagBucket{
			Desc:   "secret for user account that has admin priv",
			EnvVar: clientSecret,
		},
		configDir: flagBucket{
			Desc:   "config dir.  Default is config",
			EnvVar: configDir,
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

func getFlag(input string) (flag string) {
	flag = strings.ToLower(strings.Replace(input, "_", "-", -1))
	return
}

func getConfigDir(c *cli.Context) (cDir string) {
	cDir = c.String(getFlag(configDir))
	if cDir == "" {
		return "config"
	}
	return
}
