package viper

import (
	"fmt"
	"strings"

	"github.com/1Panel-dev/1Panel/backend/configs"
	"github.com/1Panel-dev/1Panel/backend/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() {
	v := viper.NewWithOptions()
	v.SetConfigName("app")
	v.SetConfigType("yml")
	v.AddConfigPath("/opt/1Panel/conf")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.CONF); err != nil {
			panic(err)
		}
	})
	for _, k := range v.AllKeys() {
		value := v.GetString(k)
		if strings.HasPrefix(value, "${") && strings.Contains(value, "}") {
			itemKey := strings.ReplaceAll(value[strings.Index(value, "${"):strings.Index(value, "}")], "${", "")
			v.Set(k, strings.ReplaceAll(value, fmt.Sprintf("${%s}", itemKey), v.GetString(itemKey)))
		}
	}
	serverConfig := configs.ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	global.CONF = serverConfig
}
