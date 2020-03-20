package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"testing"
)

var remote_jdkman = "https://raw.githubusercontent.com/cloudhuang/gojdkman/master/jdkman.yaml"

func TestListFolders(t *testing.T) {

	resp, err := http.Get(remote_jdkman)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)

	var jdk RemoteJDK

	err = yaml.Unmarshal(content, &jdk)
	if err != nil {
		panic(err)
	}

	println(jdk.Versions)
}

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
