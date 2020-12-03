package config

const (
	// RoleDmp .
	RoleDmp = "dmp"
	// RoleCH .
	RoleCH = "ch"
	// RoleEnclave .
	RoleEnclave = "db"
)

type (
	// Role .
	Role = string

	// Roles .
	Roles = []Role
)

// RolesAll .
var RolesAll = Roles{RoleCH, RoleDmp, RoleEnclave}
