package jdk

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"gosdkman/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var Home, _ = os.UserHomeDir()
var SdkManPath = filepath.Join(Home, ".gosdkman")
var SdkManYaml = filepath.Join(SdkManPath, "sdkman.yaml")
var currentJdkPath = filepath.Join(SdkManPath, "current")

type RemoteJDK struct {
	Versions map[string]map[string]Version `yaml:"versions"`
}

type LocalJDK struct {
	Current  string                        `yaml:"current"`
	Versions map[string]map[string]Version `yaml:"versions"`
}

type Version struct {
	Dist       string `yaml: "dist"`
	Identifier string `yaml: "identifier"`
	File       string `yaml: "file"`
}

type NewInstallVersion struct {
	vendor     string
	version    string
	dist       string
	identifier string
	file       string
}

/*
Add the new install jdk version to sdkman.yaml file
*/
func InstallNewVersion(identifier string) error {
	nv := selectAvailableJDK(identifier)

	localJDKFile, err := utils.DownloadFile(nv.file)
	if err != nil {
		return errors.New("\nFailed to download the JDK file")
	}

	err = configSDKManYaml(nv, localJDKFile, identifier)
	if err != nil {
		return errors.New("\nFailed to config the JDK path")
	}

	return UseJDKVersion(identifier)
}

func UseJDKVersion(identifier string) error {
	localJDKFile := getLocalJDKFilename(identifier)

	if localJDKFile != "" && isJDKFileExists(localJDKFile) {
		nv := selectAvailableJDK(identifier)

		err := clearOrCreateCurrentPath()
		if err != nil {
			return errors.New("\nFailed to config the JDK path")
		}

		newJDKPath := unzipJDKVersion(localJDKFile)

		// set the system environments
		err = utils.SetEnv("JAVA_HOME", filepath.Join(currentJdkPath, newJDKPath))
		if err != nil {
			return errors.New("\nFailed to config the JDK path")
		}
		err = utils.SetEnv("classpath", ".;%JAVA_HOME%\\lib")
		if err != nil {
			return errors.New("\nFailed to config the JDK path")
		}

		// update the sdkman.yaml configuration file
		err = configSDKManYaml(nv, localJDKFile, identifier)
		if err != nil {
			return err
		}

	} else {
		fmt.Println("The JDK version is not installed, will download it and configure as the current JDK version")
		return InstallNewVersion(identifier)
	}

	return nil
}

func configSDKManYaml(nv *NewInstallVersion, localJDKFile string, identifier string) error {
	newVersion := &Version{
		Identifier: nv.identifier,
		Dist:       nv.dist,
		File:       localJDKFile,
	}

	var jdk LocalJDK
	yamlFile, err := ioutil.ReadFile(SdkManYaml)
	if err != nil {
		jdk = LocalJDK{
			Current: nv.identifier,
		}
	} else {
		err = yaml.Unmarshal(yamlFile, &jdk)
		if err != nil {
			return errors.New("\nFailed to config the JDK path")
		}
	}

	if jdk.Versions == nil {
		jdk.Versions = make(map[string]map[string]Version)
	}

	if jdk.Versions[nv.vendor] == nil {
		jdk.Versions[nv.vendor] = make(map[string]Version)
	}
	jdk.Versions[nv.vendor][nv.version] = *newVersion
	jdk.Current = identifier

	marshal, _ := yaml.Marshal(jdk)

	err = ioutil.WriteFile(SdkManYaml, marshal, 0755)
	if err != nil {
		return errors.New("\nFailed to config the JDK path")
	}
	return nil
}

func unzipJDKVersion(filename string) string {
	// unzip the jdk zip file to Current folder
	utils.Unzip(filepath.Join(SdkManPath, filename), currentJdkPath)

	// list the folder
	files, err := ioutil.ReadDir(currentJdkPath)

	if err != nil {
		log.Fatal(err)
	}

	var newJDKPath string
	for _, file := range files {
		newJDKPath = file.Name()
	}

	return newJDKPath
}

func clearOrCreateCurrentPath() error {
	var err error
	if utils.Exists(currentJdkPath) {
		// delete the Current folder if exists
		err = os.RemoveAll(currentJdkPath)
	}
	if !utils.Exists(currentJdkPath) {
		err = os.Mkdir(currentJdkPath, os.FileMode(0755))
	}
	return err
}

/*
Uninstall the jdk version
*/
func UninstallVersion(identifier string) error {
	yamlFile, err := ioutil.ReadFile(SdkManYaml)
	if err != nil {
		return err
	}

	var jdk LocalJDK

	err = yaml.Unmarshal(yamlFile, &jdk)
	check(err)

	if jdk.Current == identifier {
		return errors.New(fmt.Sprintf("The JDK version '%s' currently in use, cannot uninstall", identifier))
	}

	var file = getLocalJDKFilename(identifier)

	if file != "" && isJDKFileExists(file) {
		err := os.RemoveAll(filepath.Join(SdkManPath, file))
		if err != nil {
			return errors.New("failed to delete the JDK version")
		}
	}

	marshal, _ := yaml.Marshal(jdk)

	err = ioutil.WriteFile(SdkManYaml, marshal, 0755)
	if err != nil {
		return errors.New("failed to delete the JDK version")
	}

	return nil
}

/*
List all the available JDK versions
*/
func ListAvailableJDKVersion() {
	var remoteJDK = remoteJDK()
	var localJDK = localJDK()

	fmt.Println("=====================================================================================")
	fmt.Println("Available Java Versions")
	fmt.Println("=====================================================================================")

	fmt.Printf(" %-10s | %-4s| %-12s | %-10s | %-10s | %-20s\n", "Vendor", "Use", "Version", "Dist", "Status", "Identifier")
	fmt.Println("-------------------------------------------------------------------------------------")

	for vendor, versions := range remoteJDK.Versions {
		var idx = 0
		fmt.Printf(" %-10s | ", vendor)
		for k, v := range versions {

			isUseResult := isJDKUsed(v.Identifier, localJDK.Current)
			status := isStatusInstalled(v.Identifier, localJDK)

			if idx == 0 {
				fmt.Printf("%-4s| %-12s | %-10s | %-10s | %-20s\n", isUseResult, k, v.Dist, status, v.Identifier)
			} else {
				fmt.Printf(" %-10s | %-4s| %-12s | %-10s | %-10s | %-20s\n", "", isUseResult, k, v.Dist, status, v.Identifier)
			}
			idx++
		}
	}

	fmt.Println("=====================================================================================")
	fmt.Println("Use the Identifier for installation:")
	fmt.Println("\t jdk -i 11.0.6.10.1-amaz")
}

func selectAvailableJDK(identifier string) *NewInstallVersion {
	var nv NewInstallVersion
	remoteJDK := remoteJDK()

	for vendor, versions := range remoteJDK.Versions {
		for k, v := range versions {
			if v.Identifier == identifier {
				nv = NewInstallVersion{
					vendor:     vendor,
					version:    k,
					dist:       v.Dist,
					identifier: v.Identifier,
					file:       v.File,
				}
			}

		}
	}

	return &nv
}

func remoteJDK() *RemoteJDK {
	var remoteSDKManYaml = "http://q7oqvkeuv.bkt.clouddn.com/sdkman.yaml"

	resp, err := http.Get(remoteSDKManYaml)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	var jdk RemoteJDK
	err = yaml.Unmarshal(content, &jdk)
	check(err)

	return &jdk
}

/*
Get all local JDK version
*/
func localJDK() *LocalJDK {
	yamlFile, err := ioutil.ReadFile(SdkManYaml)
	if err != nil {
		return &LocalJDK{}
	}

	var jdk LocalJDK

	err = yaml.Unmarshal(yamlFile, &jdk)
	check(err)

	return &jdk
}

func isStatusInstalled(identifier string, jdk *LocalJDK) string {
	for _, versions := range jdk.Versions {
		for _, v := range versions {
			if identifier == v.Identifier && isJDKFileExists(v.File) {
				return "Installed"
			}
		}
	}

	return ""
}

func getLocalJDKFilename(identifier string) string {
	localJDK := localJDK()

	for _, versions := range localJDK.Versions {
		for _, v := range versions {
			if identifier == v.Identifier {
				return v.File
			}
		}
	}
	return ""
}

func isJDKFileExists(file string) bool {
	return utils.Exists(filepath.Join(SdkManPath, file))
}

/*
List all the installed local LocalJDK versions
*/
func ListInstalledVersion() {
	yamlFile, err := ioutil.ReadFile(SdkManYaml)
	if err != nil {
		return
	}

	var jdk LocalJDK

	err = yaml.Unmarshal(yamlFile, &jdk)
	check(err)

	fmt.Println("=====================================================================================")
	fmt.Println("Available Java Versions")
	fmt.Println("=====================================================================================")

	fmt.Printf(" %-10s | %-4s| %-12s | %-10s | %-10s | %-20s\n", "Vendor", "Use", "Version", "Dist", "Status", "Identifier")
	fmt.Println("-------------------------------------------------------------------------------------")

	for vendor, versions := range jdk.Versions {
		var idx = 0
		fmt.Printf(" %-10s | ", vendor)
		for k, v := range versions {

			isUseResult := isJDKUsed(v.Identifier, jdk.Current)
			status := "Installed"

			if idx == 0 {
				fmt.Printf("%-4s| %-12s | %-10s | %-10s | %-20s\n", isUseResult, k, v.Dist, status, v.Identifier)
			} else {
				fmt.Printf(" %-10s | %-4s| %-12s | %-10s | %-10s | %-20s\n", "", isUseResult, k, v.Dist, status, v.Identifier)
			}
			idx++
		}
	}

	fmt.Println("=====================================================================================")
	fmt.Println("Use the Identifier for installation:")
	fmt.Println("\t jdk -i 11.0.6.10.1-amaz")
}

func isJDKUsed(identifier string, current string) string {
	if current == identifier {
		return ">>>"
	} else {
		return ""
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
