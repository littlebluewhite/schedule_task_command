package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

// loadConfig 讀取 config
func loadConfig[T any](myPath string, fileName string, configType CType) (Config T) {
	v := viper.New()
	v.AddConfigPath(myPath)
	v.SetConfigName(fileName)
	v.SetConfigType(configType.String())

	// 若有同名環境變量則使用環境變量
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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

func NewConfig[T Config](
	path, fileName string, configType CType,
) (cfg T) {
	cfg = loadConfig[T](path, fileName, configType)
	return
}

type CType int

const (
	TypeNone CType = iota
	Yaml
	Ini
)

func (hm *CType) String() string {
	return [...]string{
		"None", "yaml", "ini",
	}[*hm]
}
