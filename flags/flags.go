package flags

import (
	"os"

	"URL-Shortner/constants"

	flag "github.com/spf13/pflag"
)

var (
	env            = flag.String(constants.EnvKey, constants.EnvDefaultValue, constants.EnvUsage)
	port           = flag.Int(constants.PortKey, constants.PortDefaultValue, constants.PortUsage)
	baseConfigPath = flag.String(constants.BaseConfigPathKey, constants.BaseConfigPathDefaultValue,
		constants.BaseConfigPathUsage)
)

func init() {
	flag.Parse()
}

func Env() string {
	return *env
}

func Port() int {
	return *port
}

func BaseConfigPath() string {
	return *baseConfigPath
}

func DatabaseHost() string {
	return os.Getenv(constants.DatabaseHostKey)
}

func DatabasePassword() string {
	return os.Getenv(constants.DatabasePasswordKey)
}

func DatabaseUser() string {
	return os.Getenv(constants.DatabaseUserKey)
}

func DatabaseName() string {
	return os.Getenv(constants.DatabaseNameKey)
}
