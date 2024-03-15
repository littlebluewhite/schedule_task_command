package config

import "time"

type Config struct {
	Redis    RedisConfig    `mapstructure:"redis"`
	SQL      SQLConfig      `mapstructure:"SQL"`
	Influxdb InfluxdbConfig `mapstructure:"influxdb"`
	TestSQL  SQLConfig      `mapstructure:"testSQL"`
	Server   ServerConfig   `mapstructure:"server"`
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
	Interval     time.Duration `mapstructure:"interval"`
	CleanTime    time.Duration `mapstructure:"clean_time"`
}

type RedisConfig struct {
	Host      string `mapstructure:"host"`
	Port      string `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DB        string `mapstructure:"db"`
	IsCluster bool   `mapstructure:"is_cluster"`
}

type InfluxdbConfig struct {
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
	Org    string `mapstructure:"org"`
	Token  string `mapstructure:"token"`
	Bucket string `mapstructure:"bucket"`
}
