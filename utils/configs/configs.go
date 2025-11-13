package configs

import (
	"context"

	"URL-Shortner/constants"
	"URL-Shortner/flags"
	"URL-Shortner/log"
	"URL-Shortner/models"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config models.Config

func Get() models.Config {
	return config
}

func InitConfigs(ctx context.Context) {
	viper.SetConfigType(constants.YamlConfigType)
	viper.AddConfigPath(flags.BaseConfigPath())
	viper.SetConfigName(flags.Env() + constants.YamlFileExtension)
	err := viper.ReadInConfig()
	if err != nil {
		log.Sugar.Fatalf("Error reading config file: %v", err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.Unmarshal(&config)
		if err != nil {
			log.Sugar.Errorf("Error unmarshalling config after change: %v", err)
		} else {
			log.Sugar.Infof("Config file changed: %s", e.Name)
		}
	})
	viper.WatchConfig()

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Sugar.Fatalf("Error unmarshalling config: %v", err)
	}

	log.Sugar.Infof("Loaded configuration")
}
