package constants

const (
	PortKey          = "port"
	PortDefaultValue = 8080
	PortUsage        = "application port"

	BaseConfigPathKey          = "base-config-path"
	BaseConfigPathDefaultValue = "resources/configs"
	BaseConfigPathUsage        = "path to folder that stores your configurations"

	EnvKey          = "ENV"
	EnvDefaultValue = "dev"
	EnvUsage        = "application environment (dev, prod, test)"

	DatabaseHostKey     = "DATABASE_HOST"
	DatabasePasswordKey = "DATABASE_PASSWORD"
	DatabaseUserKey     = "DATABASE_USER"
	DatabaseNameKey     = "DATABASE_NAME"
)
