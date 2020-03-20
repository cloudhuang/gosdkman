package jdk

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var home, _ = os.UserHomeDir()
var jdkman_path = filepath.Join(home, ".jdkman")
var filename = filepath.Join(jdkman_path, "jdkman.yaml")

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

/*
Add the new install jdk version to jdkman.yaml file
*/
func InstallNewVersion(vendor string, version string, file string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	check(err)

	var jdk LocalJDK

	err = yaml.Unmarshal(yamlFile, &jdk)
	check(err)

	new_version := &Version{
		Identifier: vendor + "-" + version,
		File:       file,
	}
	jdk.Versions[vendor][version] = *new_version

	marshal, _ := yaml.Marshal(jdk)

	err = ioutil.WriteFile(filename, marshal, 0755)
	check(err)

	return nil
}

/*
Uninstall the jdk version
*/
func UninstallVersion(vendor, version string) error {
	yamlFile, err := ioutil.ReadFile(filename)
	check(err)

	var jdk LocalJDK

	err = yaml.Unmarshal(yamlFile, &jdk)
	check(err)

	versions := jdk.Versions[vendor]
	delete(versions, version)

	marshal, _ := yaml.Marshal(jdk)

	err = ioutil.WriteFile(filename, marshal, 0755)
	check(err)

	return nil
}

/*
List all the installed local LocalJDK versions
*/
func ListInstalledVersion() {
	yamlFile, err := ioutil.ReadFile(filename)
	check(err)

	var jdk LocalJDK

	err = yaml.Unmarshal(yamlFile, &jdk)
	check(err)

	fmt.Println("=====================================================================================")
	fmt.Println("Available Java Versions")
	fmt.Println("=====================================================================================")

	fmt.Printf(" %-10s | %-4s| %-12s | %-10s | %-10s | %-20s\n", "Vendor", "Use", "Version", "Dist", "Status", "Identifier")
	fmt.Println("-------------------------------------------------------------------------------------")
	// fmt.Printf("Description: %#v\n", jdk.Versions)
	//fmt.Printf("\nCurrent Version: %s\n", jdk.Current)

	for vendor, versions := range jdk.Versions {
		var idx = 0
		fmt.Printf(" %-10s | ", vendor)
		for k, v := range versions {

			isUseResult := isJDKUsed(v.Identifier, &jdk)
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

func isJDKUsed(identifier string, jdk *LocalJDK) string {
	if jdk.Current == identifier {
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
