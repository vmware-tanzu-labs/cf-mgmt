package ldap

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	l "github.com/go-ldap/ldap"
	"github.com/xchapter7x/lo"
)

var (
	userRegexp          = regexp.MustCompile(",[A-Z]+=")
	unescapeFilterRegex = regexp.MustCompile(`\\([\da-fA-F]{2}|[()\\*])`) // only match \[)*\] or \xx x=a-fA-F
)

func ParseUserCN(userDN string) (string, error) {
	dn, err := l.ParseDN(userDN)
	if err != nil {
		indexes := userRegexp.FindStringIndex(strings.ToUpper(userDN))
		if len(indexes) == 0 {
			return "", fmt.Errorf("cannot find CN for DN: %s", userDN)
		}
		cnTemp := UnescapeFilterValue(userDN[:indexes[0]])
		lo.G.Debug("CN unescaped:", cnTemp)

		escapedCN := EscapeFilterValue(cnTemp)
		lo.G.Debug("CN escaped:", escapedCN)
		return escapedCN, nil
	} else {
		attributeName := dn.RDNs[0].Attributes[0].Type
		cn := dn.RDNs[0].Attributes[0].Value
		return fmt.Sprintf("%s=%s", attributeName, cn), nil
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
