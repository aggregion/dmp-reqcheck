package config

import (
	"github.com/spf13/viper"
)

type (
	// CommonSettings .
	CommonSettings struct {
		CasTestTarget string
		CasProdTarget string
		ProxyTarget   string
		EOSTestTarget string
		EOSProdTarget string

		GatherConcurrency     int `validate:"required,min=1,max=16"`
		DefaultClickhousePort int `validate:"required,min=1,max=65535"`
	}
)

func commonSettingsValidateAndGet(v *viper.Viper) *CommonSettings {
	v.SetDefault("common.chport", 8123)
	v.SetDefault("common.concurrency", 3)
	v.SetDefault("common.proxy", "185.175.44.42:80")
	v.SetDefault("common.castest", "185.175.44.42:18765")
	v.SetDefault("common.casprod", "185.175.44.40:18765")
	v.SetDefault("common.eostest", "185.137.232.118:9999")
	v.SetDefault("common.eosprod", "185.137.232.118:8888")

	var conf = &CommonSettings{
		CasTestTarget: v.GetString("common.castest"),
		CasProdTarget: v.GetString("common.casprod"),
		ProxyTarget:   v.GetString("common.proxy"),
		EOSTestTarget: v.GetString("common.eostest"),
		EOSProdTarget: v.GetString("common.eosprod"),

		DefaultClickhousePort: v.GetInt("common.chport"),
		GatherConcurrency:     v.GetInt("common.concurrency"),
	}

	return conf
}
