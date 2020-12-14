package reports

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/zcalusic/sysinfo"
)

const versionRegExp = `([\d]+[^\d]+[\d]+)`

func parseVersionMinorMajor(version string) (major int64, minor int64) {
	versionParts := regexp.MustCompile(`\d+`).FindAllString(version, -1)
	if len(versionParts) > 0 {
		fmt.Sscanf(versionParts[0], "%d", &major)
	}
	if len(versionParts) > 1 {
		fmt.Sscanf(versionParts[1], "%d", &minor)
	}
	return
}

func getOutput(ctx context.Context, cmd string, args ...string) (string, error) {
	if cmd == "" {
		return "", errors.New("empty command")
	}

	which := exec.CommandContext(ctx, cmd, args...)
	output, err := which.CombinedOutput()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}
	return "", err
}

func getOutputFromAndTo(ctx context.Context, regExpFrom, regExpTo string, lastFindTo bool, cmd string, args ...string) (string, error) {
	output, err := getOutput(ctx, cmd, args...)
	if err != nil {
		return "", err
	}

	fromIndex := 0
	toIndex := len(output)

	if len(regExpFrom) > 0 {
		fromExp := regexp.MustCompile(regExpFrom)
		loc := fromExp.FindStringIndex(output)
		if len(loc) > 0 {
			fromIndex = loc[1]
		}
	}

	if len(regExpTo) > 0 {
		toExp := regexp.MustCompile(regExpTo)
		loc := toExp.FindAllStringIndex(output, -1)
		if len(loc) > 0 {
			if lastFindTo {
				toIndex = loc[len(loc)-1][0]
			} else {
				toIndex = loc[0][0]
			}
		}
	}

	return output[fromIndex:toIndex], nil
}

func getOutputAndRegexpFind(ctx context.Context, regExp string, cmd string, args ...string) (string, error) {
	output, err := getOutput(ctx, cmd, args...)
	if err != nil {
		return "", err
	}

	return regexp.MustCompile(regExp).FindString(output), nil
}

func getOutputAndRegexpFindAll(ctx context.Context, regExp string, cmd string, args ...string) ([]string, error) {
	output, err := getOutput(ctx, cmd, args...)
	if err != nil {
		return []string{}, err
	}

	return regexp.MustCompile(regExp).FindAllString(output, 0), nil
}

var linuxFoundAppCache = sync.Map{}

func findLinuxApps(ctx context.Context, paths []string, apps []string) (result map[string]string) {
	result = make(map[string]string)
	isNeedToFind := false
	findArgs := []string{"", "-executable", "-type", "f"}

	for inx, app := range apps {
		if valRaw, ok := linuxFoundAppCache.Load(app); ok {
			result[app] = valRaw.(string)
			continue
		}

		output, err := exec.LookPath(app)
		if err == nil {
			result[app] = string(output)
			linuxFoundAppCache.Store(app, result[app])
			continue
		}

		if inx > 0 {
			findArgs = append(findArgs, "-or")
		}
		findArgs = append(findArgs, "-iname", app)
		isNeedToFind = true
	}

	if isNeedToFind {
		for _, path := range paths {
			if len(result) == len(apps) {
				break
			}
			findArgs[0] = path
			output, err := getOutput(ctx, "find", findArgs...)
			if err == nil {
				for _, line := range strings.Split(string(output), "\n") {
					for _, app := range apps {
						if strings.HasSuffix(line, string(os.PathSeparator)+app) {
							result[app] = line
							linuxFoundAppCache.Store(app, result[app])
						}
					}
				}
			}
		}
	}

	return
}

var sysinfoOnce = sync.Once{}
var sysinfoValue *sysinfo.SysInfo

func getSysInfo() *sysinfo.SysInfo {
	sysinfoOnce.Do(func() {
		sysinfoValue = new(sysinfo.SysInfo)
		sysinfoValue.GetSysInfo()
	})

	return sysinfoValue
}
