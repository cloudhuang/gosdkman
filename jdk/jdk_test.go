package jdk

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"gosdkman/utils"
	"testing"
	"time"
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
	err := UninstallVersion("11.0.6.10.1-amaz")
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

func TestProgressingBar(t *testing.T) {
	uiprogress.Start()            // start rendering
	bar := uiprogress.AddBar(100) // Add a new bar

	// optionally, append and prepend completion and elapsed time
	bar.AppendCompleted()
	bar.PrependElapsed()

	for bar.Incr() {
		time.Sleep(time.Millisecond * 20)
	}
}
