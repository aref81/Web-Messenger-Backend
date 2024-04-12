package configs

import "github.com/spf13/viper"

type Config struct {
	Database DBConfig     `mapstructure:"database"`
	Server   ServerConfig `mapstructure:"server"`
	S3       s3Config     `mapstructure:"s3"`
}

type ServerConfig struct {
	Port      string `mapstructure:"port"`
	Address   string `mapstructure:"address"`
	SecretKey string `mapstructure:"tokenKey"`
}

type DBConfig struct {
	Port     int    `mapstructure:"port"`
	Addr     string `mapstructure:"address"`
	DBName   string `mapstructure:"dbname"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type s3Config struct {
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	Region    string `mapstructure:"region"`
	Bucket    string `mapstructure:"bucket"`
	Endpoint  string `mapstructure:"endpoint"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("configs")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return &Config{}, err
	}

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		return &Config{}, err
	}

	return &conf, nil
}
