package jdk

import (
	"testing"
)

func TestListInstalledVersion(t *testing.T) {
	ListInstalledVersion()
}

func TestAddNewVersion(t *testing.T) {

	ListInstalledVersion()

	var err error
	// Add new version
	err = InstallNewVersion("zulu", "11", "TESTPATH")
	err = InstallNewVersion("zulu", "12", "TESTPATH")
	check(err)

	err = UninstallVersion("zulu", "12")
	check(err)

	ListInstalledVersion()
}

func TestUninstallVersion(t *testing.T) {
	err := UninstallVersion("zulu", "11")
	check(err)

	ListInstalledVersion()
}