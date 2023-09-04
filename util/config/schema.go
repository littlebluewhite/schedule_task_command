package config

import "time"

type DBConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DB       string `mapstructure:"DB_DB"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"PORT"`
	ReadTimeout  time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT"`
}

type RedisConfig struct {
	Host     string `mapstructure:"REDIS_HOST"`
	Port     string `mapstructure:"REDIS_PORT"`
	User     string `mapstructure:"REDIS_USER"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       string `mapstructure:"REDIS_DB"`
}

type InfluxdbConfig struct {
	Host   string `mapstructure:"INFLUXDB_HOST"`
	Port   string `mapstructure:"INFLUXDB_PORT"`
	Org    string `mapstructure:"INFLUXDB_ORG"`
	Bucket string `mapstructure:"INFLUXDB_BUCKET"`
	Token  string `mapstructure:"INFLUXDB_TOKEN"`
}
