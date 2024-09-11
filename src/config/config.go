package config

import (
	"fmt"
	"time"
)

var ac AppConfig

type AppConfig struct {
	Env                 string
	PostgresConfig      PostgresConfig
	MessageSendDuration time.Duration
	RedisConfig         RedisConfig `yaml:"redis"`
	MessageSendClient   string
}

type PostgresConfig struct {
	Host                  string
	Port                  string
	UserName              string
	Password              string
	DBName                string
	MaxConnections        string
	MaxConnectionIdleTime string
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

var cfgs = map[string]AppConfig{
	"qa": {
		Env: "qa",
		PostgresConfig: PostgresConfig{
			Host:                  "localhost",
			Port:                  "5432",
			UserName:              "postgres",
			Password:              "postgres",
			DBName:                "messages",
			MaxConnections:        "10",
			MaxConnectionIdleTime: "30s",
		},
		MessageSendDuration: 2,
		RedisConfig: RedisConfig{
			Addr:     "localhost:6379",
			Password: "12345",
			DB:       0,
		},
		MessageSendClient: "https://webhook.site/ca7d73f7-e12b-414a-85c1-153e4f1e3914",
	},
	"prod": {
		Env: "prod",
		PostgresConfig: PostgresConfig{
			Host:                  "localhost",
			Port:                  "5432",
			UserName:              "postgres",
			Password:              "postgres",
			DBName:                "messages",
			MaxConnections:        "10",
			MaxConnectionIdleTime: "30s",
		},
		MessageSendDuration: 2,
		RedisConfig: RedisConfig{
			Addr:     "localhost:6379",
			Password: "12345",
			DB:       0,
		},
		MessageSendClient: "https://webhook.site/ca7d73f7-e12b-414a-85c1-153e4f1e3914",
	},
}

func Config(env string) AppConfig {
	cfg, exist := cfgs[env]
	if !exist {
		panic(fmt.Sprintf("config '%s' not found", env))
	}
	return cfg
}

func SetConfig(cfg AppConfig) {
	ac = cfg
}

func GetConfigs() *AppConfig {
	return &ac
}
