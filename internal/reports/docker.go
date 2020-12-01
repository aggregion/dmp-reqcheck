package reports

import (
	"context"
	"os"
	"time"
)

const (
	// DockerClientVersionStrAttr .
	DockerClientVersionStrAttr = "client_version"
	// DockerServerVersionStrAttr .
	DockerServerVersionStrAttr = "server_version"
	// DockerComposeVersionStrAttr .
	DockerComposeVersionStrAttr = "compose_version"
)

type (
	// DockerReport .
	DockerReport struct {
		ClientVersion  string `attr:"client_version"`
		ServerVersion  string `attr:"server_version"`
		ComposeVersion string `attr:"compose_version"`
	}
)

const versionRegExp = `([\d]+\.[\d]+\.[\d]+)`

func (dr *DockerReport) gatherLinux(ctx context.Context) error {
	homeDir, err := os.UserHomeDir()
	paths := findLinuxApps(
		context.Background(),
		[]string{"/usr", "/opt", homeDir},
		[]string{
			"docker",
			"dockerd",
			"docker-compose",
		},
	)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	clientVersion, err := getOutputAndRegexpFind(ctx, versionRegExp, paths["docker"], "--version")
	if err == nil {
		dr.ClientVersion = clientVersion
	}

	serverVersion, err := getOutputAndRegexpFind(ctx, versionRegExp, paths["dockerd"], "--version")
	if err == nil {
		dr.ServerVersion = serverVersion
	}

	composeVersion, err := getOutputAndRegexpFind(ctx, versionRegExp, paths["docker-compose"], "--version")
	if err == nil {
		dr.ComposeVersion = composeVersion
	}

	return nil
}

// Gather .
func (dr *DockerReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *DockerReport) GetInt64(attrName string) int64 {
	return 0
}

// GetString .
func (dr *DockerReport) GetString(attrName string) string {
	switch attrName {
	case DockerClientVersionStrAttr:
		return dr.ClientVersion
	case DockerServerVersionStrAttr:
		return dr.ServerVersion
	case DockerComposeVersionStrAttr:
		return dr.ComposeVersion
	default:
	}

	return ""
}
