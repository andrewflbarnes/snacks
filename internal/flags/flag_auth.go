package flags

import (
	"encoding/base64"
	"flag"
)

func InitAuthFlags(flagSet *flag.FlagSet) AuthFlags {
	return authFlagsImpl{
		bearer: flagSet.String("bearer", "", "The token to use for bearer auth"),
		basic:  flagSet.String("basic", "", "The user credentials to use for basic auth in the form user:pass"),
	}
}

type AuthFlags interface {
	IsAuth() bool
	GetAuth() (string, bool)
}

type authFlagsImpl struct {
	bearer *string
	basic  *string
}

func (l authFlagsImpl) IsAuth() bool {
	return l.isBasic() || l.isBearer()
}

func (l authFlagsImpl) isBasic() bool {
	return len(*l.basic) > 0
}

func (l authFlagsImpl) isBearer() bool {
	return len(*l.bearer) > 0
}

func (l authFlagsImpl) GetAuth() (string, bool) {
	if l.isBasic() {
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(*l.basic)), true
	} else if l.isBearer() {
		return "Bearer " + *l.bearer, true
	}

	return "", false
}
