package jdk

import (
	"fmt"
	"gosdkman/utils"
	"testing"
)

func TestListInstalledVersion(t *testing.T) {
	ListInstalledVersion()
}

func TestAddNewVersion(t *testing.T) {
	var err error
	// Add new version
	err = InstallNewVersion("14.27.1-zulu")
	check(err)

	ListInstalledVersion()
}

func TestUnZip(t *testing.T) {
	jdkVersion := unzipJDKVersion("zulu14.27.1-ca-jdk14-win_x64.zip")
	fmt.Println(jdkVersion)
}


func TestUninstallVersion(t *testing.T) {
	err := UninstallVersion("14.27.1-zulu")
	if err != nil {
		fmt.Println(err)
	}
}

func TestListAvailableJDKVersion(t *testing.T) {
	ListAvailableJDKVersion()
}

func TestDownloadRemoteJDKFile(t *testing.T) {
	//var identifier = "14.27.1-zulu"
	var jdkFile = "https://cdn.azul.com/zulu/bin/zulu14.27.1-ca-jdk14-win_x64.zip"
	utils.DownloadFile(jdkFile)
}
