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
	err = InstallNewVersion("zulu", "11", "amaz", "11.0.6.10.1-amaz", "TESTPATH")
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

func TestListAvailableJDKVersion(t *testing.T) {
	ListAvailableJDKVersion()
}
