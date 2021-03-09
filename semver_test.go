package main

import (
	"testing"

	"github.com/blang/semver/v4"
)

func TestParseVersionUpgrade(t *testing.T) {
	a := "bump sbt from 1.4.6 to 1.4.8"
	b := "Update discipline-scalatest from 0.0.1 -> 2.1.2"
	c := "9.4.3 (from 3.2.1)"

	avu, err := parseVersionUpgrade(a)
	if err != nil {
		t.Error(err)
	}

	bvu, err := parseVersionUpgrade(b)
	if err != nil {
		t.Error(err)
	}

	cvu, err := parseVersionUpgrade(c)
	if err != nil {
		t.Error(err)
	}

	expectEqual(avu.From, semver.MustParse("1.4.6"), t)
	expectEqual(avu.To, semver.MustParse("1.4.8"), t)

	expectEqual(bvu.From, semver.MustParse("0.0.1"), t)
	expectEqual(bvu.To, semver.MustParse("2.1.2"), t)

	expectEqual(cvu.From, semver.MustParse("3.2.1"), t)
	expectEqual(cvu.To, semver.MustParse("9.4.3"), t)
}

func TestUpgradeType(t *testing.T) {
	p := versionUpgrade{
		From: semver.MustParse("0.0.1"),
		To:   semver.MustParse("0.0.2"),
	}
	min := versionUpgrade{
		From: semver.MustParse("0.0.1"),
		To:   semver.MustParse("0.1.2"),
	}
	maj := versionUpgrade{
		From: semver.MustParse("0.0.1"),
		To:   semver.MustParse("1.1.2"),
	}

	if p.UpgradeType() != patch {
		t.Fatal("patch not detected")
	}

	if min.UpgradeType() != minor {
		t.Fatal("minor not detected")
	}

	if maj.UpgradeType() != major {
		t.Fatal("major not detected")
	}
}

func expectEqual(a, b semver.Version, t *testing.T) {
	if a.String() != b.String() {
		t.Fatalf("%v != %v", a, b)
	}
}
