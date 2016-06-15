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
		CreateCommand("create-orgs", runCreateOrgs, defaultFlags(), eh),
		CreateCommand("create-spaces", runCreateSpaces, defaultFlags(), eh),
		CreateCommand("update-spaces", runUpdateSpaces, defaultFlags(), eh),
	}

	return app
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
	if theSystemDomain, theUserID, thePassword, theConfigDir, theSecret, theError := getRequiredFields(c); theError == nil {
		uaaManager := uaa.NewDefaultUAAManager(theSystemDomain, theUserID)
		cfToken := uaaManager.GetCFToken(thePassword)
		uaacToken := uaaManager.GetUAACToken(theSecret)
		orgManager := organization.NewManager(theSystemDomain, cfToken, uaacToken)
		err = orgManager.CreateOrgs(theConfigDir)
	} else {
		err = theError
	}
	return
}

func runCreateSpaces(c *cli.Context) (err error) {
	if theSystemDomain, theUserID, thePassword, theConfigDir, theSecret, theError := getRequiredFields(c); theError == nil {
		uaaManager := uaa.NewDefaultUAAManager(theSystemDomain, theUserID)
		cfToken := uaaManager.GetCFToken(thePassword)
		uaacToken := uaaManager.GetUAACToken(theSecret)
		orgManager := space.NewManager(theSystemDomain, cfToken, uaacToken)
		err = orgManager.CreateSpaces(theConfigDir)
	} else {
		err = theError
	}
	return
}

func runUpdateSpaces(c *cli.Context) (err error) {
	if theSystemDomain, theUserID, thePassword, theConfigDir, theSecret, theError := getRequiredFields(c); theError == nil {
		uaaManager := uaa.NewDefaultUAAManager(theSystemDomain, theUserID)
		cfToken := uaaManager.GetCFToken(thePassword)
		uaacToken := uaaManager.GetUAACToken(theSecret)
		orgManager := space.NewManager(theSystemDomain, cfToken, uaacToken)
		err = orgManager.UpdateSpaces(theConfigDir)
	} else {
		err = theError
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

func getRequiredFields(c *cli.Context) (sysDomain, user, pwd, config, secret string, err error) {
	sysDomain = c.String(getFlag(systemDomain))
	user = c.String(getFlag(userID))
	pwd = c.String(getFlag(password))
	config = c.String(getFlag(configDir))
	secret = c.String(getFlag(clientSecret))

	if sysDomain == "" ||
		user == "" ||
		pwd == "" ||
		config == "" ||
		secret == "" {
		err = fmt.Errorf("Must set system-domain, user-id, password, config-dir, client-secret properties")
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
