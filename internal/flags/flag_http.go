package flags

import (
	"encoding/base64"
	"flag"
	"regexp"
	"strings"
)

var (
	headerRe = regexp.MustCompile(" *: *")
)

func InitHttpFlags(flagSet *flag.FlagSet) HttpFlags {
	return httpFlagsImpl{
		bearer:  flagSet.String("bearer", "", "The token to use for bearer auth"),
		basic:   flagSet.String("basic", "", "The user credentials to use for basic auth in the form user:pass"),
		headers: flagSet.String("headers", "", "Additional HTTP headers to repeat for the attack, separate with a double pipe i.e. \"||\""),
	}
}

type HttpFlags interface {
	IsAuth() bool
	GetAuth() (string, bool)
	GetHeaders() map[string]string
}

type httpFlagsImpl struct {
	bearer  *string
	basic   *string
	headers *string
}

func (l httpFlagsImpl) IsAuth() bool {
	return l.isBasic() || l.isBearer()
}

func (l httpFlagsImpl) isBasic() bool {
	return len(*l.basic) > 0
}

func (l httpFlagsImpl) isBearer() bool {
	return len(*l.bearer) > 0
}

func (l httpFlagsImpl) GetAuth() (string, bool) {
	if l.isBasic() {
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(*l.basic)), true
	} else if l.isBearer() {
		return "Bearer " + *l.bearer, true
	}

	return "", false
}

func (l httpFlagsImpl) GetHeaders() map[string]string {
	headers := make(map[string]string)

	for _, header := range strings.Split(*l.headers, "||") {
		if len(header) == 0 {
			continue
		}
		hkv := headerRe.Split(header, -1)
		headers[hkv[0]] = hkv[1]
	}

	if l.isBasic() {
		headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(*l.basic))
	} else if l.isBearer() {
		headers["Authorization"] = "Bearer " + *l.bearer
	}

	return headers
}
