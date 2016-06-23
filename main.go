package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cf-mgmt/generated"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
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
}

//InitializeManager -
func InitializeManager(c *cli.Context) (cfMgmt CFMgmt, err error) {
	sysDomain := c.String(getFlag(systemDomain))
	user := c.String(getFlag(userID))
	pwd := c.String(getFlag(password))
	config := c.String(getFlag(configDir))
	secret := c.String(getFlag(clientSecret))

	if sysDomain == "" ||
		user == "" ||
		pwd == "" ||
		config == "" ||
		secret == "" {
		err = fmt.Errorf("Must set system-domain, user-id, password, config-dir, client-secret properties")
	} else {
		cfMgmt = CFMgmt{}
		cfMgmt.UAAManager = uaa.NewDefaultUAAManager(sysDomain, user)
		cfToken := cfMgmt.UAAManager.GetCFToken(pwd)
		uaacToken := cfMgmt.UAAManager.GetUAACToken(secret)
		cfMgmt.OrgManager = organization.NewManager(sysDomain, cfToken, uaacToken)
		cfMgmt.SpaceManager = space.NewManager(sysDomain, cfToken, uaacToken)
		cfMgmt.ConfigDir = config
	}

	return
}

const (
	systemDomain string = "SYSTEM_DOMAIN"
	userID       string = "USER_ID"
	password     string = "PASSWORD"
	clientSecret string = "CLIENT_SECRET"
	configDir    string = "CONFIG_DIR"
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
		CreateGeneratePipelineCommand(runGeneratePipeline, eh),
		CreateCommand("create-orgs", runCreateOrgs, defaultFlags(), eh),
		CreateCommand("create-spaces", runCreateSpaces, defaultFlags(), eh),
		CreateCommand("update-spaces", runUpdateSpaces, defaultFlags(), eh),
		CreateCommand("update-space-users", runUpdateSpaceUsers, defaultFlags(), eh),
		CreateCommand("update-org-users", runUpdateOrgUsers, defaultFlags(), eh),
	}

	return app
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
	fmt.Println("fly -t lite set-pipeline -p cf-mgmt -c pipeline.yml --load-vars-from=vars.yml --var \"git_private_key=$(cat ~/.ssh/id_rsa)\"")
	return
}

func createFile(assetName, fileName string) (err error) {
	var f *os.File
	var fileBytes []byte
	if fileBytes, err = generated.Asset(fmt.Sprintf("pipeline/%s", assetName)); err == nil {
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
	var cfMgmt CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.OrgManager.CreateOrgs(cfMgmt.ConfigDir)
	}
	return
}

func runCreateSpaces(c *cli.Context) (err error) {
	var cfMgmt CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.CreateSpaces(cfMgmt.ConfigDir)
	}
	return
}

func runUpdateSpaces(c *cli.Context) (err error) {
	var cfMgmt CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.UpdateSpaces(cfMgmt.ConfigDir)
	}
	return
}

func runUpdateSpaceUsers(c *cli.Context) (err error) {
	var cfMgmt CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.SpaceManager.UpdateSpaceUsers(cfMgmt.ConfigDir)
	}
	return
}

func runUpdateOrgUsers(c *cli.Context) (err error) {
	var cfMgmt CFMgmt
	if cfMgmt, err = InitializeManager(c); err == nil {
		err = cfMgmt.OrgManager.UpdateOrgUsers(cfMgmt.ConfigDir)
	}
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
			Desc:   "config dir.  Default is .",
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
