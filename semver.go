package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/blang/semver/v4"
)

type versionUpgrade struct {
	From semver.Version
	To   semver.Version
}

type upgradeType int

const (
	noChange = iota
	patch
	minor
	major
)

func (vu *versionUpgrade) UpgradeType() upgradeType {
	if vu.To.Major > vu.From.Major {
		return major
	}

	if vu.To.Minor > vu.From.Minor {
		return minor
	}

	if vu.To.Patch > vu.From.Patch {
		return patch
	}

	return noChange
}

func parseUpgradeType(s string) (upgradeType, error) {
	parsed, ok := map[string]upgradeType{
		"major": major,
		"minor": minor,
		"patch": patch,
	}[strings.TrimSpace(strings.ToLower(s))]

	if !ok {
		return noChange, fmt.Errorf("unrecognised allowed update value %v", s)
	}

	return parsed, nil
}

// extract two versions from a string, order them and return a VersionUpgrade
func parseVersionUpgrade(s string) (*versionUpgrade, error) {
	semVerRegex := regexp.MustCompile(`([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?`)
	matches := semVerRegex.FindAllString(s, 2)

	if len(matches) != 2 {
		return nil, fmt.Errorf("unable to parse 2 versions from %v", s)
	}

	a := semver.MustParse(matches[0])
	b := semver.MustParse(matches[1])

	if a.LTE(b) {
		return &versionUpgrade{a, b}, nil
	}

	return &versionUpgrade{b, a}, nil
}

func allowed(allowed, proposed upgradeType) bool {
	return proposed <= allowed
}
