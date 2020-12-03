package config

import (
	"fmt"
	"strings"

	"github.com/aggregion/dmp-reqcheck/pkg/utils"
	"github.com/spf13/viper"
)

type (
	// HostSettings .
	HostSettings struct {
		Roles []string `validate:"required"`

		DefaultClickhousePort int

		Hosts map[string]string

		IsListen bool
	}
)

func hostSettingsValidateAndGet(v *viper.Viper) *HostSettings {
	v.SetDefault("defaults.chport", 8123)

	var conf = &HostSettings{
		Roles: strings.Split(v.GetString("host.roles"), ","),
		Hosts: map[string]string{},

		DefaultClickhousePort: v.GetInt("defaults.chport"),

		IsListen: v.GetBool("host.listen"),
	}

	for _, host := range strings.Split(v.GetString("host.hosts"), ",") {
		hostParts := strings.Split(host, ":")
		if len(hostParts) == 1 && hostParts[0] == "" {
			continue
		}

		if len(hostParts) != 2 {
			panic(fmt.Sprintf("the host argument [%s] is not valid, not specified a role", host))
		}
		conf.Hosts[hostParts[0]] = hostParts[1]

		if !utils.IsIntersectStrs(RolesAll, Roles{hostParts[0]}) {
			panic(fmt.Sprintf("the host argument [%s] is not valid, unknown host role (supported %s)", host, RolesAll))
		}
	}

	// check roles
	if utils.IsIntersectStrs(conf.Roles, Roles{RoleDmp, RoleEnclave}) {
		if len(conf.Hosts[RoleCH]) == 0 {
			panic("for role dmp you should specify --hosts ch:host")
		}
		if len(conf.Hosts[RoleEnclave]) == 0 {
			panic("for role dmp you should specify --hosts ch:enclave")
		}
	}

	utils.MustValidate(conf)

	return conf
}
