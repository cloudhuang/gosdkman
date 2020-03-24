package main

import (
	"fmt"
	"github.com/devfacet/gocmd"
	jdk "gosdkman/jdk"
	utils "gosdkman/utils"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//var identifier = "jdk-13-openjdk.zip"
var identifier = "openjdk-14_windows-x64_bin.zip"

var home, _ = os.UserHomeDir()
var sdkManPath = filepath.Join(home, ".sdkman")
var currentJdkPath = filepath.Join(sdkManPath, "current")

func main() {
	flags := struct {
		Help      bool   `short:"h" long:"help" description:"Display usage" global:"true"`
		Version   bool   `short:"v" long:"version" description:"Display version"`
		Install   string `short:"i" long:"install" description:"Install the new JDK version"`
		Uninstall string `short:"d" long:"uninstall" description:"Uninstall the JDK version"`
		Use       string `short:"u" long:"use" description:"Use this JDK version"`
		List      bool   `short:"l" long:"list" description:"List all the available versions"`
	}{}

	// List command
	gocmd.HandleFlag("List", func(cmd *gocmd.Cmd, args []string) error {
		jdk.ListAvailableJDKVersion()
		return nil
	})

	gocmd.HandleFlag("Use", func(cmd *gocmd.Cmd, args []string) error {
		identifier = args[0]
		fmt.Println(identifier)
		return nil
	})

	gocmd.HandleFlag("Install", func(cmd *gocmd.Cmd, args []string) error {
		identifier = args[0]
		err := jdk.InstallNewVersion(identifier)
		if err == nil {
			fmt.Println("The new JDK version installed success, please restart the console.")
		} else {

		}
		return nil
	})

	gocmd.HandleFlag("Uninstall", func(cmd *gocmd.Cmd, args []string) error {
		identifier = args[0]
		fmt.Println(identifier)
		return nil
	})

	// Init the app
	gocmd.New(gocmd.Options{
		Name:        "GoSDKMan",
		Version:     "1.0.0",
		Description: "GoSDKMan is a command-line tool which allows you to easily install, manage, and work with multiple Java environments for windows.",
		Flags:       &flags,
		ConfigType:  gocmd.ConfigTypeAuto,
	})

	//installNewJDK()

	//printJavaVersion()
}

func installNewJDK() {
	err := os.MkdirAll(sdkManPath, os.ModePerm)
	if err != nil {
		fmt.Println("Create sdkman folder failed.")
	}

	printJavaVersion()

	// delete the Current folder is exists
	err = os.RemoveAll(currentJdkPath)

	// unzip the jdk zip file to Current folder
	utils.Unzip(filepath.Join(sdkManPath, identifier), currentJdkPath)

	// list the folder
	files, err := ioutil.ReadDir(currentJdkPath)

	if err != nil {
		log.Fatal(err)
	}

	var new_jdk_path string
	for _, file := range files {
		new_jdk_path = file.Name()
	}

	jdk_path := filepath.Join(currentJdkPath, new_jdk_path)
	fmt.Printf("\nThe new JDK HOME is %s", jdk_path)

	err_path := utils.SetEnv("JAVA_HOME", jdk_path)
	err_classpath := utils.SetEnv("classpath", ".;%JAVA_HOME%\\lib")
	if err_path != nil && err_classpath != nil {
		fmt.Errorf("failed to config JDK in system: %v - %v", err_path, err_classpath)
	}
}

func printJavaVersion() {
	var_java_home := os.Getenv("JAVA_HOME")
	fmt.Printf("\nThe Current JAVA_HOME: %s", var_java_home)

	var_classpath := os.Getenv("classpath")
	fmt.Printf("\nThe Current classpath: %s", var_classpath)
}
