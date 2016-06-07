package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/pivotalservices/cf-mgmt/organization"
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
	configFile   string = "CONFIG_FILE"
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
		CreateSyncOrgsCommand(eh),
	}

	return app
}

//CreateSyncOrgsCommand -
func CreateSyncOrgsCommand(eh *ErrorHandler) (command cli.Command) {
	desc := fmt.Sprintf("sync-orgs")
	command = cli.Command{
		Name:        "sync-orgs",
		Usage:       "sync orgs with what is defined in config",
		Description: desc,
		Action:      runSyncOrgs,
		Flags:       syncOrgsFlags(),
	}
	return
}

func syncOrgsFlags() (flags []cli.Flag) {
	var flagList = map[string]flagBucket{
		systemDomain: flagBucket{
			Desc:   "system domain",
			EnvVar: systemDomain,
		},
		userID: flagBucket{
			Desc:   "user id that can create/delete orgs",
			EnvVar: userID,
		},
		password: flagBucket{
			Desc:   "password for user account that can create/delete orgs",
			EnvVar: password,
		},
		configFile: flagBucket{
			Desc:   "config file for orgs.  Default is orgs.yml in current directory",
			EnvVar: configFile,
		},
	}

	flags = buildFlags(flagList)
	return
}

func runSyncOrgs(c *cli.Context) (err error) {
	var token string
	var theConfigFile = "orgs.yml"
	theSytemDomain := c.String(getFlag(systemDomain))
	theUserID := c.String(getFlag(userID))
	thePassword := c.String(getFlag(password))
	if c.IsSet(getFlag(configFile)) {
		theConfigFile = c.String(getFlag(configFile))
	}
	uaamanager := uaa.NewDefaultUAAManager(theSytemDomain, theUserID, thePassword)
	if token, err = uaamanager.GetToken(); err == nil {
		orgManager := organization.NewDefaultOrgManager(theSytemDomain, token)
		err = orgManager.SyncOrgs(theConfigFile)
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
