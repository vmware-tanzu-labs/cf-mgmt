package ldap

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	l "github.com/go-ldap/ldap/v3"
	"github.com/xchapter7x/lo"
)

var (
	userRegexp          = regexp.MustCompile(",[A-Z]+=")
	unescapeFilterRegex = regexp.MustCompile(`\\([\da-fA-F]{2}|[()\\*])`) // only match \[)*\] or \xx x=a-fA-F
)

func getUserAttributeName(userDN string) string {
	parts := strings.Split(userDN, "=")
	return parts[0]
}

func ParseUserCN(userDN string) (string, string, error) {
	dn, err := l.ParseDN(userDN)
	if err != nil {
		indexes := userRegexp.FindStringIndex(strings.ToUpper(userDN))
		if len(indexes) == 0 {
			return "", "", fmt.Errorf("cannot find CN for DN: %s", userDN)
		}
		index := indexes[0]
		cnTemp := UnescapeFilterValue(userDN[:index])
		lo.G.Debug("CN unescaped:", cnTemp)

		escapedCN := EscapeFilterValue(cnTemp)
		lo.G.Debug("CN escaped:", escapedCN)
		return escapedCN, userDN[index+1:], nil
	} else {
		userAttributeName := getUserAttributeName(userDN)
		searchBase := ""
		attributeName := dn.RDNs[0].Attributes[0].Type
		cn := dn.RDNs[0].Attributes[0].Value
		for _, rdn := range dn.RDNs {
			attrType := rdn.Attributes[0].Type
			if !strings.EqualFold(userAttributeName, attrType) {
				if len(searchBase) > 0 {
					searchBase = searchBase + ","
				}
				searchBase = searchBase + fmt.Sprintf("%s=%s", attrType, rdn.Attributes[0].Value)
			}
		}
		return fmt.Sprintf("%s=%s", attributeName, cn), searchBase, nil
	}
}

func EscapeFilterValue(filter string) string {
	return l.EscapeFilter(strings.Replace(filter, "\\", "", -1))
}

func UnescapeFilterValue(filter string) string {
	repl := unescapeFilterRegex.ReplaceAllFunc(
		[]byte(filter),
		func(match []byte) []byte {
			// \( \) \\ \*
			if len(match) == 2 {
				return []byte{match[1]}
			}
			// had issues with Decode, TODO fix to use Decode?.
			res, _ := hex.DecodeString(string(match[1:]))
			return res
		},
	)
	return string(repl)
}
