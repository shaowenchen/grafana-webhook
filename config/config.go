package config

import "github.com/spf13/viper"

var Config = Options{}

func ReadConfig() {
	Config = Options{
		Gin: GinOption{viper.GetString("gin.runmode")},
	}
}

type Options struct {
	Gin GinOption
}

type GinOption struct {
	RunMode string
}
