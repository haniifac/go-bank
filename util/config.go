package util

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

// Stores all configuration of the application using viper
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	// experimental feature viper 1.20.1 for binding automaticEnv to viper env registry
	// see for more details: https://github.com/spf13/viper/issues/1797
	v := viper.NewWithOptions(viper.ExperimentalBindStruct())

	// read from system env and set to viper,
	// but without set it to viper env registry (e.g. `viper.GetString("DB_DRIVER")` will not work)
	v.AutomaticEnv()

	if os.Getenv("IS_PRODUCTION") != "true" {
		v.AddConfigPath(path)
		v.SetConfigName("app")
		v.SetConfigType("env")

		if err = v.ReadInConfig(); err != nil {
			return
		}
	}

	err = v.Unmarshal(&config)
	return
}
