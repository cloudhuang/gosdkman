package jdk

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var home, _ = os.UserHomeDir()
var jdkmanPath = filepath.Join(home, ".jdkman")
var jdkmanYaml = filepath.Join(jdkmanPath, "jdkman.yaml")

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
Add the new install jdk version to jdkman.yaml file
*/
func InstallNewVersion(identifier string) error {

	nv := selectAvailableJDK(identifier)

	var jdk LocalJDK
	yamlFile, err := ioutil.ReadFile(jdkmanYaml)
	if err != nil {
		jdk = LocalJDK{
			Current: nv.identifier,
		}
	} else {
		err = yaml.Unmarshal(yamlFile, &jdk)
		check(err)
		jdk.Current = identifier
	}

	newVersion := &Version{
		Identifier: nv.identifier,
		Dist:       nv.dist,
		File:       nv.file,
	}

	if jdk.Versions == nil {
		jdk.Versions = make(map[string]map[string]Version)
	}

	if jdk.Versions[nv.vendor] == nil {
		jdk.Versions[nv.vendor] = make(map[string]Version)
	}
	jdk.Versions[nv.vendor][nv.version] = *newVersion

	marshal, _ := yaml.Marshal(jdk)

	err = ioutil.WriteFile(jdkmanYaml, marshal, 0755)
	check(err)

	// TODO Download the remote JDK file
	return nil
}

/*
Uninstall the jdk version
*/
func UninstallVersion(identifier string) error {
	yamlFile, err := ioutil.ReadFile(jdkmanYaml)
	if err != nil {
		return err
	}

	var jdk LocalJDK

	err = yaml.Unmarshal(yamlFile, &jdk)
	check(err)

	var file string
	for vendor, versions := range jdk.Versions {
		for k, v := range versions {
			if v.Identifier == identifier {
				delete(jdk.Versions[vendor], k)
				file = v.File
			}
		}
	}

	if file != "" {
		// TODO delete the jdk file
		fmt.Printf("the jdk file: %s\n", file)
	}

	marshal, _ := yaml.Marshal(jdk)

	err = ioutil.WriteFile(jdkmanYaml, marshal, 0755)
	check(err)

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
			status := isInstalled(v.Identifier, localJDK)

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
	fmt.Println("\t remoteJDK -i 11.0.6.10.1-amaz")
}

func selectAvailableJDK(identifier string) *NewInstallVersion {
	var nv NewInstallVersion
	remoteJDK := remoteJDK()

	for vendor, versions := range remoteJDK.Versions {
		fmt.Printf(" %-10s | ", vendor)
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

	// TODO download the remote jdk file
}

func remoteJDK() *RemoteJDK {
	var remoteJDKManYaml = "https://raw.githubusercontent.com/cloudhuang/gojdkman/master/jdkman.yaml"

	resp, err := http.Get(remoteJDKManYaml)
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
	yamlFile, err := ioutil.ReadFile(jdkmanYaml)
	if err != nil {
		return &LocalJDK{}
	}

	var jdk LocalJDK

	err = yaml.Unmarshal(yamlFile, &jdk)
	check(err)

	return &jdk
}

func isInstalled(identifier string, jdk *LocalJDK) string {
	for _, versions := range jdk.Versions {
		for _, v := range versions {
			if identifier == v.Identifier {
				return "Installed"
			}
		}
	}

	return ""
}

/*
List all the installed local LocalJDK versions
*/
func ListInstalledVersion() {
	yamlFile, err := ioutil.ReadFile(jdkmanYaml)
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
