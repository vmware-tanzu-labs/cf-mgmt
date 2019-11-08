package ldap

import (
	"encoding/hex"
	"fmt"
	"strings"

	l "github.com/go-ldap/ldap"
	"github.com/pkg/errors"
)

func ParseUserCN(userDN string) (string, error) {
	dn, err := l.ParseDN(userDN)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to parse userDN %s", userDN)
	}
	attributeName := dn.RDNs[0].Attributes[0].Type
	cn := dn.RDNs[0].Attributes[0].Value
	return fmt.Sprintf("%s=%s", attributeName, cn), nil
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
