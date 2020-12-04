package config

import (
	"log"
	"strings"

	"github.com/aggregion/dmp-reqcheck/pkg/utils"
	"github.com/spf13/viper"
)

type (
	// HostSettings .
	HostSettings struct {
		Roles []string `validate:"required,min=1,dive,oneof=ch dmp enclave"`

		GatherConcurrency     int `validate:"required,min=1,max=16"`
		DefaultClickhousePort int `validate:"required,min=1,max=65535"`

		Hosts map[string]string

		IsListen bool
	}
)

func hostSettingsValidateAndGet(v *viper.Viper, isListenContext bool) *HostSettings {
	v.SetDefault("defaults.chport", 8123)
	v.SetDefault("defaults.concurrency", 3)

	var conf = &HostSettings{
		Roles: utils.FilterEmptyStrs(strings.Split(v.GetString("host.roles"), ",")),
		Hosts: map[string]string{},

		DefaultClickhousePort: v.GetInt("defaults.chport"),

		IsListen:          v.GetBool("host.listen"),
		GatherConcurrency: v.GetInt("defaults.concurrency"),
	}

	if !isListenContext {
		for _, host := range utils.FilterEmptyStrs(strings.Split(v.GetString("host.hosts"), ",")) {
			hostParts := strings.Split(host, ":")
			if len(hostParts) != 2 {
				log.Fatalf("the host argument [%s] is not valid, not specified a role", host)
			}
			conf.Hosts[hostParts[0]] = hostParts[1]

			if !utils.IsIntersectStrs(RolesAll, Roles{hostParts[0]}) {
				log.Fatalf("the host argument [%s] is not valid, unknown host role (supported %s)", host, RolesAll)
			}
		}

		// check roles
		if utils.IsIntersectStrs(conf.Roles, Roles{RoleDmp, RoleEnclave}) {
			if len(conf.Hosts[RoleCH]) == 0 {
				log.Fatalf("for roles dmp or enclave you should specify --hosts ch:host")
			}
		}
		if utils.IsIntersectStrs(conf.Roles, Roles{RoleEnclave}) {
			if len(conf.Hosts[RoleDmp]) == 0 {
				log.Fatalf("for role enclave you should specify --hosts dmp:host")
			}
		}
		if utils.IsIntersectStrs(conf.Roles, Roles{RoleDmp}) {
			if len(conf.Hosts[RoleEnclave]) == 0 {
				log.Fatalf("for role dmp you should specify --hosts enclave:host")
			}
		}
	}

	utils.MustValidate(conf)

	return conf
}
