package configs

import (
	"time"

	"github.com/spf13/viper"
)

type Conf struct {
	IpMaxReqPerSecond    uint              `mapstructure:"IP_MAX_REQ_PER_SECOND"`
	IpBlockDuration      time.Duration     `mapstructure:"IP_BLOCK_DURATION"`
	TokenMaxReqPerSecond uint              `mapstructure:"TOKEN_MAX_REQ_PER_SECOND"`
	TokenBlockDuration   time.Duration     `mapstructure:"TOKEN_BLOCK_DURATION"`
	ServerPort           string            `mapstructure:"SERVER_PORT"`
	RedisURL             string            `mapstructure:"REDIS_URL"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	viper.AutomaticEnv()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
