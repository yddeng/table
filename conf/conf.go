package conf

import "github.com/BurntSushi/toml"

type Config struct {
	WSAddr   string
	HttpAddr string
	LoadDir  string
	DBConfig *DBConfig
}

type DBConfig struct {
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbDataBase string
}

var config *Config

func LoadConfig(path string) {
	config = &Config{}
	_, err := toml.DecodeFile(path, config)
	if err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return config
}
