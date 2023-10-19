package config

import (
	"github.com/spf13/viper"
	"log"
	"path"
)

// loadConfig 讀取.env環境變數檔
func loadConfig[T any](myPath string, fileName string) (Config T) {
	v := viper.New()
	v.AddConfigPath(myPath)
	v.SetConfigName(fileName)
	v.SetConfigType("yaml")

	// 若有同名環境變量則使用環境變量
	v.AutomaticEnv()
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("can not load config: " + err.Error())
	}
	err = v.Unmarshal(&Config)
	if err != nil {
		log.Fatal("can not load config: " + err.Error())
	}
	return
}

func NewConfig[T DBConfig | ServerConfig | RedisConfig | InfluxdbConfig](rootPath, folder, fileName string) (cfg T) {
	cfg = loadConfig[T](path.Join(rootPath, folder), fileName)
	return
}
