package config

import "time"

type Config struct {
	Conn   ConnConfig   `mapstructure:"conn"`
	Server ServerConfig `mapstructure:"server"`
}

type ConnConfig struct {
	Redis    RedisConfig    `mapstructure:"redis"`
	SQL      SQLConfig      `mapstructure:"SQL"`
	Influxdb InfluxdbConfig `mapstructure:"influxdb"`
	TestSQL  SQLConfig      `mapstructure:"testSQL"`
}

type SQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       string `mapstructure:"db"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type RedisConfig struct {
	Host      string   `mapstructure:"host"`
	Ports     []string `mapstructure:"port"`
	User      string   `mapstructure:"user"`
	Password  string   `mapstructure:"password"`
	DB        string   `mapstructure:"db"`
	IsCluster bool     `mapstructure:"is_cluster"`
}

type InfluxdbConfig struct {
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
	Org    string `mapstructure:"org"`
	Token  string `mapstructure:"token"`
	Bucket string `mapstructure:"bucket"`
}
