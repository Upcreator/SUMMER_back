package initializers

import (
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_NAME"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	JwtSecret      string `mapstructure:"JWT_SECRET"`
	Stand          string `mapstructure:"STAND"`
	FrontendUrl    string `mapstructure:"FRONTEND_URL"`
	CookieDomain   string `mapstructure:"COOKIE_DOMAIN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	AppConfig = config
	return config, err
}
