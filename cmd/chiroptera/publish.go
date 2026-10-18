package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Linkernel0x/chiroptera/pkg/pack"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(publishCmd)

	publishCmd.Flags().StringP("file", "f", "./chiro/package.toml", "Path to the package configuration file")
	publishCmd.Flags().IntP("port", "p", -1, "Port to listen on if self hosting")
}

type publishOptions struct {
	Path    string
	Hosting bool
	Port    int
}

var publishCmd = &cobra.Command{
	Use:   "publish [dev@package#version]",
	Short: "Publish a recipe ",
	Run: func(cmd *cobra.Command, args []string) {
		hasNoArgs := len(args) == 0
		hasNoFlags := cmd.Flags().NFlag() == 0

		if hasNoArgs && hasNoFlags {
			fmt.Println("No argument or flag found, loading wizard...")
			runPublishWizard()
			return
		}

		path, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Println("Error while obtaining 'file' flag: ", err)
			return
		}

		port, err := cmd.Flags().GetInt("port")
		var hosting bool
		if err != nil {
			fmt.Println("Error while obtaining 'port' flag: ", err)
			return
		}
		if port >= 1024 && port <= 65535 {
			hosting = true
		} else {
			hosting = false
			port = -1
		}

		options := publishOptions{
			Path:    path,
			Hosting: hosting,
			Port:    port,
		}

		runPublish(options)
	},
}

func runPublishWizard() {
	questions := []*survey.Question{
		{
			Name: "path",
			Prompt: &survey.Input{
				Renderer: survey.Renderer{},
				Message:  "The wizard will look for a recipe at ./.chiro/package.toml, if the path of your recipe is different write it here. (you should be in the root directory of your project)",
				Default:  "./.chiro/package.toml",
			},
			Validate: survey.Required,
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Renderer: survey.Renderer{},
				Message:  "Will your package be self hosted on this machine? If so write down the port.\n This will open an http server on this machine at the selected door, if you don't know what that means DO NOT ENABLE this feature. It's up to you to protect the port behind proxies and firewalls.",
				Default:  "",
			},
		},
		{
			Name: "confirm",
			Prompt: &survey.Confirm{
				Renderer: survey.Renderer{},
				Message:  "Do you want to proceed?",
			},
		},
	}

	answers := struct {
		Path    string
		Port    string
		Confirm bool
	}{}

	err := survey.Ask(questions, &answers)
	if err != nil {
		fmt.Println("Wizard quit because of an error: ", err)
		return
	}

	if answers.Confirm {
		fmt.Printf("\n✨ Excellent! I'll proceed to publish your package!\n\n")
	} else {
		return
	}

	var path string
	if answers.Path == "" || (!strings.HasPrefix(answers.Path, "/") && !strings.HasPrefix(answers.Path, ".") && !strings.HasPrefix(answers.Path, "~")) {
		path = "./.chiro/chiro.toml"
	} else {
		path = answers.Path
	}

	var hosting bool
	port, err := strconv.Atoi(answers.Port)
	if err != nil && answers.Port != "" {
		fmt.Println("Port must be an integer (1024-65535)")
		return
	}
	if port >= 1024 && port <= 65535 && answers.Port != "" {
		hosting = true
	} else {
		hosting = false
		port = -1
	}

	options := publishOptions{
		Path:    path,
		Hosting: hosting,
		Port:    port,
	}

	runPublish(options)
}

func runPublish(options publishOptions) {
	_, err := recipe.ParseRecipe(options.Path)
	if err != nil {
		fmt.Println("Error while parsing recipe: ", err)
		return
	}
}
