package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type organizations struct {
	AdminPassword    string `mapstructure:"admin_password"`
	AdminEmail       string `mapstructure:"admin_email"`
	AdminUsername    string `mapstructure:"admin_username"`
	AdminDisplayname string `mapstructure:"admin_displayname"`
	OrgName          string `mapstructure:"org_name"`
	OrgDescription   string `mapstructure:"org_description"`
}

type Bootstrap struct {
	Initial organizations `mapstructure:"initial"`
}

type JWT struct {
	SigningKey string `mapstructure:"signing_key"`
	Duration   string `mapstructure:"duration"`
}

type Config struct {
	DatabaseURL string    `mapstructure:"database_url"`
	Port        string    `mapstructure:"server_port"`
	WithProxy   bool      `mapstructure:"with_proxy"`
	Bootstrap   Bootstrap `mapstructure:"bootstrap"`
	UserJWT     JWT       `mapstructure:"user_jwt"`
}

func Load(configpaths ...string) (*Config, error) {
	v := viper.New()

	v.SetEnvPrefix("cmms")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetConfigName("app")
	v.SetConfigType("yaml")

	SetDefaults(v)

	for _, path := range configpaths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("Failed to read the configuration file: %s", err)
	}

	config := &Config{}
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// SetDefaults sets the default values to empty fields
func SetDefaults(v *viper.Viper) {
	v.SetDefault("database_url", "")
	v.SetDefault("server_port", ":5600")
	v.SetDefault("with_proxy", false)
	v.SetDefault("bootstrap.initial.admin_email", "")
	v.SetDefault("bootstrap.initial.admin_username", "")
	v.SetDefault("bootstrap.initial.admin_displayname", "")
	v.SetDefault("bootstrap.initial.admin_password", "")
	v.SetDefault("bootstrap.initial.org_name", "")
	v.SetDefault("bootstrap.initial.org_displayname", "")
	v.SetDefault("user_jwt.signing_key", "")
	v.SetDefault("user_jwt.duration", "720h")
}
