package kernel

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Raw   string
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Version) Compare(other Version) int {
	if v.Major != other.Major {
		return v.Major - other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor - other.Minor
	}
	return v.Patch - other.Patch
}

func (v Version) GTE(other Version) bool { return v.Compare(other) >= 0 }
func (v Version) LTE(other Version) bool { return v.Compare(other) <= 0 }

func ParseVersion(s string) (Version, error) {
	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	m := re.FindStringSubmatch(s)
	if len(m) < 4 {
		return Version{}, fmt.Errorf("cannot parse kernel version from: %s", s)
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	patch, _ := strconv.Atoi(m[3])
	return Version{Major: major, Minor: minor, Patch: patch, Raw: s}, nil
}

func GetVersion() (Version, error) {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return Version{}, fmt.Errorf("failed to read /proc/version: %w", err)
	}
	return ParseVersion(strings.TrimSpace(string(data)))
}

func GetArch() string {
	return runtime.GOARCH
}

func IsContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	data, err := os.ReadFile("/proc/1/cgroup")
	if err == nil {
		s := string(data)
		if strings.Contains(s, "docker") || strings.Contains(s, "lxc") || strings.Contains(s, "kubepods") {
			return true
		}
	}
	return false
}
