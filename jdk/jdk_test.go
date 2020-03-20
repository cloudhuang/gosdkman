package jdk

import (
	"fmt"
	"testing"
)

func TestListInstalledVersion(t *testing.T) {
	ListInstalledVersion()
}

func TestAddNewVersion(t *testing.T) {
	var err error
	// Add new version
	err = InstallNewVersion("11.0.6.10.1-amaz")
	check(err)

	ListInstalledVersion()
}

func TestUninstallVersion(t *testing.T) {
	err := UninstallVersion("11.0.6.10.1-amaz")
	if err != nil {
		fmt.Println(err)
	}
}

func TestListAvailableJDKVersion(t *testing.T) {
	ListAvailableJDKVersion()
}
