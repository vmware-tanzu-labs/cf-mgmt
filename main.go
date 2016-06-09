package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
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

const (
	systemDomain string = "SYSTEM_DOMAIN"
	userID       string = "USER_ID"
	password     string = "PASSWORD"
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
			Action: func(c *cli.Context) {
				cli.ShowVersion(c)
			},
		},
		CreateOrgsCommand(eh),
		CreateSpacesCommand(eh),
	}

	return app
}

//CreateOrgsCommand -
func CreateOrgsCommand(eh *ErrorHandler) (command cli.Command) {
	desc := fmt.Sprintf("create-orgs")
	command = cli.Command{
		Name:        "create-orgs",
		Usage:       "create orgs with what is defined in config",
		Description: desc,
		Action:      runCreateOrgs,
		Flags:       createOrgsFlags(),
	}
	return
}

func createOrgsFlags() (flags []cli.Flag) {
	var flagList = buildDefaultFlags()
	flags = buildFlags(flagList)
	return
}

func runCreateOrgs(c *cli.Context) (err error) {
	var token, theSystemDomain, theUserID, thePassword string
	var theConfigDir = "."

	if theSystemDomain, theUserID, thePassword, theConfigDir, err = getRequiredFields(c); err != nil {
		return
	}

	uaamanager := uaa.NewDefaultUAAManager(theSystemDomain, theUserID, thePassword)
	if token, err = uaamanager.GetToken(); err == nil {
		orgManager := organization.NewManager(theSystemDomain, token)
		err = orgManager.CreateOrgs(theConfigDir)
	}
	return
}

//CreateSpacesCommand -
func CreateSpacesCommand(eh *ErrorHandler) (command cli.Command) {
	desc := fmt.Sprintf("create-spaces")
	command = cli.Command{
		Name:        "create-spaces",
		Usage:       "create spaces with what is defined in config",
		Description: desc,
		Action:      runCreateSpaces,
		Flags:       createSpacesFlags(),
	}
	return
}

func createSpacesFlags() (flags []cli.Flag) {
	var flagList = buildDefaultFlags()
	flags = buildFlags(flagList)
	return
}

func runCreateSpaces(c *cli.Context) (err error) {
	var token, theSystemDomain, theUserID, thePassword string
	var theConfigDir = "."

	if theSystemDomain, theUserID, thePassword, theConfigDir, err = getRequiredFields(c); err != nil {
		return
	}

	uaamanager := uaa.NewDefaultUAAManager(theSystemDomain, theUserID, thePassword)
	if token, err = uaamanager.GetToken(); err == nil {
		orgManager := space.NewManager(theSystemDomain, token)
		err = orgManager.CreateSpaces(theConfigDir)
	}
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
		configDir: flagBucket{
			Desc:   "config dir.  Default is .",
			EnvVar: configDir,
		},
	}
	return
}

func getRequiredFields(c *cli.Context) (sysDomain, user, pwd, config string, err error) {
	sysDomain = c.String(getFlag(systemDomain))
	user = c.String(getFlag(userID))
	pwd = c.String(getFlag(password))
	if c.IsSet(getFlag(configDir)) {
		config = c.String(getFlag(configDir))
	}
	if sysDomain == "" ||
		user == "" ||
		pwd == "" {
		err = fmt.Errorf("Must set system-domain, user-id and password properties")
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
