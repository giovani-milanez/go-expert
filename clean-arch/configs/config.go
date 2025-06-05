package configs

import "github.com/spf13/viper"

type Conf struct {
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBUser            string `mapstructure:"DB_USER"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBName            string `mapstructure:"DB_NAME"`
	WebServerPort     string `mapstructure:"WEB_SERVER_PORT"`
	GRPCServerPort    string `mapstructure:"GRPC_SERVER_PORT"`
	GraphQLServerHost string `mapstructure:"GRAPHQL_SERVER_HOST"`
	GraphQLServerPort string `mapstructure:"GRAPHQL_SERVER_PORT"`
	RabbitMqUser      string `mapstructure:"RABBITMQ_USER"`
	RabbitMqPassword  string `mapstructure:"RABBITMQ_PASSWORD"`
	RabbitMqHost      string `mapstructure:"RABBITMQ_HOST"`
	RabbitMqPort      string `mapstructure:"RABBITMQ_PORT"`
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
