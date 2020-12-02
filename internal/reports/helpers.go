package reports

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/zcalusic/sysinfo"
)

const versionRegExp = `([\d]+[^\d]+[\d]+[^\d]+[\d]+)`

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

func getReportAttrFromField(v reflect.Value) interface{} {
	kind := v.Type().Kind()

	switch kind {
	case reflect.String:
		return v.String()
	case reflect.Int:
	case reflect.Int32:
	case reflect.Int64:
		return v.Int()
	}

	return nil
}

func getReportInterfaceAttr(obj interface{}, attrName string) interface{} {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()
	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		attrsPrefix := t.Field(fieldIndex).Tag.Get("attrMap")
		if len(attrsPrefix) > 0 && strings.HasPrefix(attrName, attrsPrefix) {
			mp := v.Field(fieldIndex)
			key := reflect.ValueOf(strings.TrimPrefix(attrName, attrsPrefix))
			mv := mp.MapIndex(key)
			if !mv.IsValid() {
				return nil
			}

			return getReportAttrFromField(mp.MapIndex(key).Elem())
		}
		attrs := t.Field(fieldIndex).Tag.Get("attr")
		if attrs == attrName {
			return getReportAttrFromField(v.Field((fieldIndex)))
		}
	}

	return nil
}

func getReportStrAttr(obj interface{}, attrName string) string {
	attrVal := getReportInterfaceAttr(obj, attrName)
	if val, ok := attrVal.(string); ok {
		return val
	}

	return ""
}

func getReportIntAttr(obj interface{}, attrName string) int64 {
	attrVal := getReportInterfaceAttr(obj, attrName)
	if val, ok := attrVal.(int64); ok {
		return val
	}

	return 0
}

// GetReportAttrs .
func GetReportAttrs(obj interface{}) (result map[string]interface{}) {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()
	result = make(map[string]interface{})
	for fieldIndex := 0; fieldIndex < t.NumField(); fieldIndex++ {
		attrsPrefix := t.Field(fieldIndex).Tag.Get("attrMap")
		if len(attrsPrefix) > 0 {
			mp := v.Field(fieldIndex)
			for _, key := range mp.MapKeys() {
				result[attrsPrefix+key.String()] = getReportAttrFromField(mp.MapIndex(key).Elem())
			}
			continue
		}

		attr := t.Field(fieldIndex).Tag.Get("attr")
		if len(attr) > 0 {
			result[attr] = getReportAttrFromField(v.Field(fieldIndex))
		}
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

		output, err := getOutput(ctx, "which", app)
		if err == nil {
			result[app] = strings.TrimSpace(string(output))
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
