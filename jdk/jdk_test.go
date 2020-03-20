package jdk

import (
	"testing"
)

func TestListInstalledVersion(t *testing.T) {
	ListInstalledVersion()
}

func TestAddNewVersion(t *testing.T) {
	var err error
	// Add new version
	err = InstallNewVersion("13.0.1-zulu")
	check(err)

	ListInstalledVersion()
}

func TestUninstallVersion(t *testing.T) {
	err := UninstallVersion("11.0.6.10.1-amaz")
	check(err)
}

func TestListAvailableJDKVersion(t *testing.T) {
	ListAvailableJDKVersion()
}
