package main

import (
	"fmt"
	"github.com/devfacet/gocmd"
	jdk "gosdkman/jdk"
)

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
		identifier := args[0]
		fmt.Println(identifier)
		return nil
	})

	gocmd.HandleFlag("Install", func(cmd *gocmd.Cmd, args []string) error {
		identifier := args[0]
		err := jdk.InstallNewVersion(identifier)
		if err == nil {
			fmt.Println("The new JDK version installed success, please restart the console.")
		} else {
			fmt.Println(err.Error())
		}
		return nil
	})

	gocmd.HandleFlag("Uninstall", func(cmd *gocmd.Cmd, args []string) error {
		identifier := args[0]
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
}
