package reports

import (
	"context"
	"os"
	"path"
	"time"
)

const (
	// DockerClientVersionStrAttr .
	DockerClientVersionStrAttr = "client_version"
	// DockerClientMajorVersionIntAttr .
	DockerClientMajorVersionIntAttr = "client_version_major"
	// DockerClientMinorVersionIntAttr .
	DockerClientMinorVersionIntAttr = "client_version_minor"
	// DockerServerVersionStrAttr .
	DockerServerVersionStrAttr = "server_version"
	// DockerServerMajorVersionIntAttr .
	DockerServerMajorVersionIntAttr = "server_version_major"
	// DockerServerMinorVersionIntAttr .
	DockerServerMinorVersionIntAttr = "server_version_minor"
	// DockerComposeVersionStrAttr .
	DockerComposeVersionStrAttr = "compose_version"
	// DockerComposeMajorVersionIntAttr .
	DockerComposeMajorVersionIntAttr = "compose_version_major"
	// DockerComposeMinorVersionIntAttr .
	DockerComposeMinorVersionIntAttr = "compose_version_minor"
)

type (
	// DockerReport .
	DockerReport struct {
		clientVersion      string `attr:"client_version"`
		clientVersionMajor int64  `attr:"client_version_major"`
		clientVersionMinor int64  `attr:"client_version_minor"`

		serverVersion      string `attr:"server_version"`
		serverVersionMajor int64  `attr:"server_version_major"`
		serverVersionMinor int64  `attr:"server_version_minor"`

		composeVersion      string `attr:"compose_version"`
		composeVersionMajor int64  `attr:"compose_version_major"`
		composeVersionMinor int64  `attr:"compose_version_minor"`
	}
)

func (dr *DockerReport) gatherLinux(ctx context.Context) error {
	dr.clientVersion = ""
	dr.clientVersionMajor, dr.clientVersionMinor = 0, 0
	dr.serverVersion = ""
	dr.serverVersionMajor, dr.serverVersionMinor = 0, 0
	dr.composeVersion = ""
	dr.composeVersionMajor, dr.composeVersionMinor = 0, 0

	homeDir, err := os.UserHomeDir()
	paths := findLinuxApps(
		context.Background(),
		[]string{"/usr", "/opt", path.Join(homeDir, ".local"), path.Join(homeDir, "bin")},
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
		dr.clientVersion = clientVersion
		dr.clientVersionMajor, dr.clientVersionMinor = parseVersionMinorMajor(dr.clientVersion)
	}

	serverVersion, err := getOutputAndRegexpFind(ctx, versionRegExp, paths["dockerd"], "--version")
	if err == nil {
		dr.serverVersion = serverVersion
		dr.serverVersionMajor, dr.serverVersionMinor = parseVersionMinorMajor(dr.serverVersion)
	} else if len(dr.clientVersion) > 0 {
		// TODO: try to get from docker version, server engine
	}

	composeVersion, err := getOutputAndRegexpFind(ctx, versionRegExp, paths["docker-compose"], "--version")
	if err == nil {
		dr.composeVersion = composeVersion
		dr.composeVersionMajor, dr.composeVersionMinor = parseVersionMinorMajor(dr.composeVersion)
	}

	return nil
}

// Gather .
func (dr *DockerReport) Gather(ctx context.Context) error {
	return dr.gatherLinux(ctx)
}

// GetInt64 .
func (dr *DockerReport) GetInt64(attrName string) int64 {
	return getReportIntAttr(dr, attrName)
}

// GetString .
func (dr *DockerReport) GetString(attrName string) string {
	return getReportStrAttr(dr, attrName)
}
