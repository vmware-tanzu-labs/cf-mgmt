package configcommands

import (
	"fmt"
	"strconv"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

// BaseConfigCommand - commmand that specifies config-dir
type BaseConfigCommand struct {
	ConfigDirectory string `long:"config-dir" env:"CONFIG_DIR" default:"config" description:"Name of the config directory"`
}

type UserRole struct {
	UserRoleAdd
	LDAPUsersToRemove  []string `long:"ldap-user-to-remove" description:"Ldap User to remove, specify multiple times"`
	UsersToRemove      []string `long:"user-to-remove" description:"User to remove, specify multiple times"`
	SamlUsersToRemove  []string `long:"saml-user-to-remove" description:"SAML user to remove, specify multiple times"`
	LDAPGroupsToRemove []string `long:"ldap-group-to-remove" description:"LDAP Group to remove, specify multiple times"`
	AADGroupsToRemove  []string `long:"azure-ad-group-to-remove" description:"Azure AD Group to remove, specify multiple times"`
}

type UserRoleAdd struct {
	LDAPUsers  []string `long:"ldap-user" description:"Ldap User to add, specify multiple times"`
	Users      []string `long:"user" description:"User to add, specify multiple times"`
	SamlUsers  []string `long:"saml-user" description:"SAML user to add, specify multiple times"`
	LDAPGroups []string `long:"ldap-group" description:"LDAP Group to add, specify multiple times"`
	AADGroups  []string `long:"azure-ad-group" description:"Azure AD Group to add, specify multiple times"`
}

type ServiceAccess struct {
	ServiceAccessAdd
	ServiceNameToRemove string `long:"service-to-remove" description:"*****DEPRECATED, use 'cf-mgmt-config global service-access' ***** - name of service to remove"`
}

type ServiceAccessAdd struct {
	ServiceName string   `long:"service" description:"*****DEPRECATED, use 'cf-mgmt-config global service-access' ***** - Service Name to add"`
	Plans       []string `long:"plans" description:"*****DEPRECATED, use 'cf-mgmt-config global service-access' ***** - plans to add, empty list will add all plans"`
}

type Metadata struct {
	LabelKey            []string `long:"label" description:"Label to add, can specify multiple"`
	LabelValue          []string `long:"label-value" description:"Label value to add, can specify multiple but need to match number of label args"`
	AnnotationKey       []string `long:"annotation" description:"Annotation to add, can specify multiple"`
	AnnotationValue     []string `long:"annotation-value" description:"Annotation value to add, can specify multiple but need to match number of annotation args"`
	LabelsToRemove      []string `long:"labels-to-remove" description:"name of label to remove"`
	AnnotationsToRemove []string `long:"annotations-to-remove" description:"name of annotation to remove"`
}

type OrgQuota struct {
	EnableOrgQuota             string `long:"enable-org-quota" description:"Enable the Org Quota in the config" choice:"true" choice:"false"`
	MemoryLimit                string `long:"memory-limit" description:"An Org's memory limit in Megabytes"`
	InstanceMemoryLimit        string `long:"instance-memory-limit" description:"Global Org Application instance memory limit in Megabytes"`
	TotalRoutes                string `long:"total-routes" description:"Total Routes capacity for an Org"`
	TotalServices              string `long:"total-services" description:"Total Services capacity for an Org"`
	PaidServicesAllowed        string `long:"paid-service-plans-allowed" description:"Allow paid services to appear in an org" choice:"true" choice:"false"`
	TotalPrivateDomains        string `long:"total-private-domains" description:"Total Private Domain capacity for an Org"`
	TotalReservedRoutePorts    string `long:"total-reserved-route-ports" description:"Total Reserved Route Ports capacity for an Org"`
	TotalServiceKeys           string `long:"total-service-keys" description:"Total Service Keys capacity for an Org"`
	AppInstanceLimit           string `long:"app-instance-limit" description:"App Instance Limit an Org"`
	AppTaskLimit               string `long:"app-task-limit" description:"App Task Limit an Org"`
	LogRateLimitBytesPerSecond string `long:"log-rate-limit-bytes-per-second" description:"Log Rate limit per app for an org"`
}

type NamedOrgQuota struct {
	MemoryLimit                string `long:"memory-limit" description:"An Org's memory limit in Megabytes"`
	InstanceMemoryLimit        string `long:"instance-memory-limit" description:"Global Org Application instance memory limit in Megabytes"`
	TotalRoutes                string `long:"total-routes" description:"Total Routes capacity for an Org"`
	TotalServices              string `long:"total-services" description:"Total Services capacity for an Org"`
	PaidServicesAllowed        string `long:"paid-service-plans-allowed" description:"Allow paid services to appear in an org" choice:"true" choice:"false"`
	TotalPrivateDomains        string `long:"total-private-domains" description:"Total Private Domain capacity for an Org"`
	TotalReservedRoutePorts    string `long:"total-reserved-route-ports" description:"Total Reserved Route Ports capacity for an Org"`
	TotalServiceKeys           string `long:"total-service-keys" description:"Total Service Keys capacity for an Org"`
	AppInstanceLimit           string `long:"app-instance-limit" description:"App Instance Limit an Org"`
	AppTaskLimit               string `long:"app-task-limit" description:"App Task Limit an Org"`
	LogRateLimitBytesPerSecond string `long:"log-rate-limit-bytes-per-second" description:"Log Rate limit per app for an org"`
}

type SpaceQuota struct {
	EnableSpaceQuota           string `long:"enable-space-quota" description:"Enable the Space Quota in the config" choice:"true" choice:"false"`
	MemoryLimit                string `long:"memory-limit" description:"An Space's memory limit in Megabytes"`
	InstanceMemoryLimit        string `long:"instance-memory-limit" description:"Space Application instance memory limit in Megabytes"`
	TotalRoutes                string `long:"total-routes" description:"Total Routes capacity for an Space"`
	TotalServices              string `long:"total-services" description:"Total Services capacity for an Space"`
	PaidServicesAllowed        string `long:"paid-service-plans-allowed" description:"Allow paid services to appear in an Space" choice:"true" choice:"false"`
	TotalReservedRoutePorts    string `long:"total-reserved-route-ports" description:"Total Reserved Route Ports capacity for an Space"`
	TotalServiceKeys           string `long:"total-service-keys" description:"Total Service Keys capacity for an Space"`
	AppInstanceLimit           string `long:"app-instance-limit" description:"App Instance Limit for a space"`
	AppTaskLimit               string `long:"app-task-limit" description:"App Task Limit for a space"`
	LogRateLimitBytesPerSecond string `long:"log-rate-limit-bytes-per-second" description:"Log Rate limit per app for a space"`
}

type NamedSpaceQuota struct {
	MemoryLimit                string `long:"memory-limit" description:"An Space's memory limit in Megabytes"`
	InstanceMemoryLimit        string `long:"instance-memory-limit" description:"Space Application instance memory limit in Megabytes"`
	TotalRoutes                string `long:"total-routes" description:"Total Routes capacity for an Space"`
	TotalServices              string `long:"total-services" description:"Total Services capacity for an Space"`
	PaidServicesAllowed        string `long:"paid-service-plans-allowed" description:"Allow paid services to appear in an Space" choice:"true" choice:"false"`
	TotalReservedRoutePorts    string `long:"total-reserved-route-ports" description:"Total Reserved Route Ports capacity for an Space"`
	TotalServiceKeys           string `long:"total-service-keys" description:"Total Service Keys capacity for a Space"`
	AppInstanceLimit           string `long:"app-instance-limit" description:"App Instance Limit for a space"`
	AppTaskLimit               string `long:"app-task-limit" description:"App Task Limit for a space"`
	LogRateLimitBytesPerSecond string `long:"log-rate-limit-bytes-per-second" description:"Log Rate limit per app for a space"`
}

func updateUsersBasedOnRole(userMgmt *config.UserMgmt, currentLDAPGroups []string, currentAADGroups []string, userRole *UserRole, errorString *string) {
	userMgmt.LDAPGroups = removeFromSlice(addToSlice(currentLDAPGroups, userRole.LDAPGroups, errorString), userRole.LDAPGroupsToRemove)
	userMgmt.AADGroups = removeFromSlice(addToSlice(currentAADGroups, userRole.AADGroups, errorString), userRole.AADGroupsToRemove)
	userMgmt.Users = removeFromSlice(addToSlice(userMgmt.Users, userRole.Users, errorString), userRole.UsersToRemove)
	userMgmt.SamlUsers = removeFromSlice(addToSlice(userMgmt.SamlUsers, userRole.SamlUsers, errorString), userRole.SamlUsersToRemove)
	userMgmt.LDAPUsers = removeFromSlice(addToSlice(userMgmt.LDAPUsers, userRole.LDAPUsers, errorString), userRole.LDAPUsersToRemove)
	userMgmt.LDAPGroup = ""
}

func addUsersBasedOnRole(userMgmt *config.UserMgmt, currentLDAPGroups []string, currentAADGroups []string, userRole *UserRoleAdd, errorString *string) {
	userMgmt.LDAPGroups = addToSlice(currentLDAPGroups, userRole.LDAPGroups, errorString)
	userMgmt.AADGroups = addToSlice(currentAADGroups, userRole.AADGroups, errorString)
	userMgmt.Users = addToSlice(userMgmt.Users, userRole.Users, errorString)
	userMgmt.SamlUsers = addToSlice(userMgmt.SamlUsers, userRole.SamlUsers, errorString)
	userMgmt.LDAPUsers = addToSlice(userMgmt.LDAPUsers, userRole.LDAPUsers, errorString)
	userMgmt.LDAPGroup = ""
}

func convertToInt(parameterName string, currentValue *int, proposedValue string, errorString *string) {
	if proposedValue == "" {
		return
	}
	i, err := strconv.Atoi(proposedValue)
	if err != nil {
		*errorString += fmt.Sprintf("\n--%s must be an integer instead of [%s]", parameterName, proposedValue)
		return
	}
	*currentValue = i

}

func convertToGB(parameterName string, currentValue *string, proposedValue string, defaultValue string, errorString *string) {
	if *currentValue == "" && proposedValue == "" {
		*currentValue = defaultValue
	}
	if proposedValue == "" {
		return
	}
	val, err := config.StringToMegabytes(proposedValue)
	if err != nil {
		*errorString += fmt.Sprintf("\n--%s must be an integer instead of [%s]", parameterName, proposedValue)
		return
	}
	*currentValue = val
}

func convertToFormattedInt(parameterName string, currentValue *string, proposedValue string, defaultValue string, errorString *string) {
	if *currentValue == "" && proposedValue == "" {
		*currentValue = defaultValue
	}
	if proposedValue == "" {
		return
	}
	val, err := config.ToInteger(proposedValue)
	if err != nil {
		*errorString += fmt.Sprintf("\n--%s must be an integer instead of [%s]", parameterName, proposedValue)
		return
	}
	*currentValue = config.AsString(val)
}

func convertToBool(parameterName string, currentValue *bool, proposedValue string, errorString *string) {
	if proposedValue == "" {
		return
	}
	b, err := strconv.ParseBool(proposedValue)
	if err != nil {
		*errorString += fmt.Sprintf("\n--%s must be an boolean instead of [%s]", parameterName, proposedValue)
		return
	}
	*currentValue = b
}

func addToSlice(theSlice, sliceToAdd []string, errorString *string) []string {
	checkForDuplicates(sliceToAdd, errorString)
	sliceToReturn := theSlice
	valuesThatExist := sliceToMap(theSlice)
	for _, val := range sliceToAdd {
		if _, ok := valuesThatExist[val]; !ok && val != "" {
			sliceToReturn = append(sliceToReturn, val)
		}
	}
	return sliceToReturn
}

func checkForDuplicates(slice []string, errorString *string) {
	sliceMap := make(map[string]string)
	for _, val := range slice {
		if _, ok := sliceMap[val]; ok && val != "" {
			*errorString += fmt.Sprintf("\n--value [%s] cannot be specified more than once %v", val, slice)
		} else {
			sliceMap[val] = val
		}
	}
}

func removeFromSlice(theSlice, sliceToRemove []string) []string {
	var sliceToReturn []string
	valuesToRemove := sliceToMap(sliceToRemove)
	for _, val := range theSlice {
		if _, ok := valuesToRemove[val]; !ok && val != "" {
			sliceToReturn = append(sliceToReturn, val)
		}
	}
	return sliceToReturn
}

func sliceToMap(theSlice []string) map[string]string {
	theMap := make(map[string]string)
	for _, val := range theSlice {
		theMap[val] = val
	}
	return theMap
}

func updateOrgQuotaConfig(namedQuota string, clearNamedQuota bool, orgConfig *config.OrgConfig, orgQuota OrgQuota, errorString *string) {
	if clearNamedQuota {
		orgConfig.NamedQuota = ""
		orgConfig.MemoryLimit = config.UNLIMITED
		orgConfig.InstanceMemoryLimit = config.UNLIMITED
		orgConfig.TotalRoutes = config.UNLIMITED
		orgConfig.TotalServices = config.UNLIMITED
		orgConfig.PaidServicePlansAllowed = false
		orgConfig.TotalReservedRoutePorts = config.UNLIMITED
		orgConfig.TotalServiceKeys = config.UNLIMITED
		orgConfig.AppInstanceLimit = config.UNLIMITED
		orgConfig.AppTaskLimit = config.UNLIMITED
		orgConfig.LogRateLimitBytesPerSecond = config.UNLIMITED
	}
	if namedQuota != "" {
		orgConfig.NamedQuota = namedQuota
		orgConfig.EnableOrgQuota = false
		orgConfig.MemoryLimit = ""
		orgConfig.InstanceMemoryLimit = ""
		orgConfig.TotalRoutes = ""
		orgConfig.TotalServices = ""
		orgConfig.PaidServicePlansAllowed = false
		orgConfig.TotalReservedRoutePorts = ""
		orgConfig.TotalServiceKeys = ""
		orgConfig.AppInstanceLimit = ""
		orgConfig.AppTaskLimit = ""
		orgConfig.LogRateLimitBytesPerSecond = ""
	} else {
		convertToBool("enable-org-quota", &orgConfig.EnableOrgQuota, orgQuota.EnableOrgQuota, errorString)
		convertToGB("memory-limit", &orgConfig.MemoryLimit, orgQuota.MemoryLimit, "100M", errorString)
		convertToGB("instance-memory-limit", &orgConfig.InstanceMemoryLimit, orgQuota.InstanceMemoryLimit, config.UNLIMITED, errorString)
		convertToFormattedInt("total-routes", &orgConfig.TotalRoutes, orgQuota.TotalRoutes, config.UNLIMITED, errorString)
		convertToFormattedInt("total-services", &orgConfig.TotalServices, orgQuota.TotalServices, config.UNLIMITED, errorString)
		convertToBool("paid-service-plans-allowed", &orgConfig.PaidServicePlansAllowed, orgQuota.PaidServicesAllowed, errorString)
		convertToFormattedInt("total-private-domains", &orgConfig.TotalPrivateDomains, orgQuota.TotalPrivateDomains, config.UNLIMITED, errorString)
		convertToFormattedInt("total-reserved-route-ports", &orgConfig.TotalReservedRoutePorts, orgQuota.TotalReservedRoutePorts, config.UNLIMITED, errorString)
		convertToFormattedInt("total-service-keys", &orgConfig.TotalServiceKeys, orgQuota.TotalServiceKeys, config.UNLIMITED, errorString)
		convertToFormattedInt("app-instance-limit", &orgConfig.AppInstanceLimit, orgQuota.AppInstanceLimit, config.UNLIMITED, errorString)
		convertToFormattedInt("app-task-limit", &orgConfig.AppTaskLimit, orgQuota.AppTaskLimit, config.UNLIMITED, errorString)
		convertToFormattedInt("log-rate-limit-bytes-per-second", &orgConfig.LogRateLimitBytesPerSecond, orgQuota.LogRateLimitBytesPerSecond, config.UNLIMITED, errorString)
	}
}

func updateSpaceQuotaConfig(namedQuota string, clearNamedQuota bool, spaceConfig *config.SpaceConfig, spaceQuota SpaceQuota, errorString *string) {
	if clearNamedQuota {
		spaceConfig.NamedQuota = ""
		spaceConfig.MemoryLimit = config.UNLIMITED
		spaceConfig.InstanceMemoryLimit = config.UNLIMITED
		spaceConfig.TotalRoutes = config.UNLIMITED
		spaceConfig.TotalServices = config.UNLIMITED
		spaceConfig.PaidServicePlansAllowed = false
		spaceConfig.TotalReservedRoutePorts = config.UNLIMITED
		spaceConfig.TotalServiceKeys = config.UNLIMITED
		spaceConfig.AppInstanceLimit = config.UNLIMITED
		spaceConfig.AppTaskLimit = config.UNLIMITED
		spaceConfig.LogRateLimitBytesPerSecond = config.UNLIMITED
	}
	if namedQuota != "" {
		spaceConfig.NamedQuota = namedQuota
		spaceConfig.EnableSpaceQuota = false
		spaceConfig.MemoryLimit = ""
		spaceConfig.InstanceMemoryLimit = ""
		spaceConfig.TotalRoutes = ""
		spaceConfig.TotalServices = ""
		spaceConfig.PaidServicePlansAllowed = false
		spaceConfig.TotalReservedRoutePorts = ""
		spaceConfig.TotalServiceKeys = ""
		spaceConfig.AppInstanceLimit = ""
		spaceConfig.AppTaskLimit = ""
		spaceConfig.LogRateLimitBytesPerSecond = ""
	} else {
		convertToBool("enable-space-quota", &spaceConfig.EnableSpaceQuota, spaceQuota.EnableSpaceQuota, errorString)
		convertToGB("memory-limit", &spaceConfig.MemoryLimit, spaceQuota.MemoryLimit, config.UNLIMITED, errorString)
		convertToGB("instance-memory-limit", &spaceConfig.InstanceMemoryLimit, spaceQuota.InstanceMemoryLimit, config.UNLIMITED, errorString)
		convertToFormattedInt("total-routes", &spaceConfig.TotalRoutes, spaceQuota.TotalRoutes, config.UNLIMITED, errorString)
		convertToFormattedInt("total-services", &spaceConfig.TotalServices, spaceQuota.TotalServices, config.UNLIMITED, errorString)
		convertToBool("paid-service-plans-allowed", &spaceConfig.PaidServicePlansAllowed, spaceQuota.PaidServicesAllowed, errorString)
		convertToFormattedInt("total-reserved-route-ports", &spaceConfig.TotalReservedRoutePorts, spaceQuota.TotalReservedRoutePorts, config.UNLIMITED, errorString)
		convertToFormattedInt("total-service-keys", &spaceConfig.TotalServiceKeys, spaceQuota.TotalServiceKeys, config.UNLIMITED, errorString)
		convertToFormattedInt("app-instance-limit", &spaceConfig.AppInstanceLimit, spaceQuota.AppInstanceLimit, config.UNLIMITED, errorString)
		convertToFormattedInt("log-rate-limit-bytes-per-second", &spaceConfig.LogRateLimitBytesPerSecond, spaceQuota.LogRateLimitBytesPerSecond, config.UNLIMITED, errorString)
	}
}

func validateASGsExist(configuredASGs []config.ASGConfig, asgs []string, errorString *string) {
	asgMap := make(map[string]string)
	for _, configuredASG := range configuredASGs {
		asgMap[configuredASG.Name] = configuredASG.Name
	}
	for _, asg := range asgs {
		if _, ok := asgMap[asg]; !ok {
			*errorString += fmt.Sprintf("\n--[%s.json] does not exist in asgs directory", asg)
		}
	}
}
