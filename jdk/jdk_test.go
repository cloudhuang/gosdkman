package jdk

import (
	"fmt"
	"gosdkman/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

func TestLength(t *testing.T) {
	fmt.Println(len("AdoptOpenJDK"))
}

func TestCopyFolder(t *testing.T) {
	fmt.Println(currentJdkPath)
	// list the folder
	files, err := ioutil.ReadDir(currentJdkPath)

	if err != nil {
		log.Fatal(err)
	}

	var newJDKPath string
	for _, file := range files {
		newJDKPath = file.Name()
	}

	fmt.Println(newJDKPath)

	err = utils.CopyDir(filepath.Join(currentJdkPath, newJDKPath), currentJdkPath)
	os.RemoveAll(filepath.Join(currentJdkPath, newJDKPath))
}