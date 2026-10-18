package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Linkernel0x/chiroptera/pkg/crypto"
	"github.com/Linkernel0x/chiroptera/pkg/filesystem"
	recipe "github.com/Linkernel0x/chiroptera/pkg/pack"
	"github.com/spf13/cobra"
)

func init() {
	recipeCmd.AddCommand(createCmd)
	recipeCmd.AddCommand(validateCmd)

	rootCmd.AddCommand(recipeCmd)
}

var recipeCmd = &cobra.Command{
	Use:   "recipe",
	Short: "Interact with package recipes",
	Args:  cobra.ExactArgs(1),
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new package",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🦇  Welcome to Chiroptera recipe wizard! This will help you create a package.toml file for your programs. In order to create and publish your package you should already have complete the 'chiro setup', if that's not the case please stop the wizard and run 'chiro setup'.")

		questions := []*survey.Question{
			{
				Name: "path",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "What's the absolute path to your package's root directory?",
					Default:  "e.g. my-awesome-software",
				},
				Validate: survey.Required,
			},
			{
				Name: "name",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "What's the name of your package?",
					Default:  "e.g. my-awesome-software",
				},
				Validate: survey.Required,
			},
			{
				Name: "version",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "What's the version of your software?",
					Default:  "1.0.0",
				},
				Validate: survey.Required,
			},
			{
				Name: "description",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "Do you want to add a description?",
				},
			},
			{
				Name: "url",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "[source] At what direct url will be the package available?",
				},
				Validate: survey.Required,
			},
			{
				Name: "file",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "[source] What's the path to your package (this is needed to calculate its hash)?",
				},
				Validate: survey.Required,
			},
			{
				Name: "torrent",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "[source] Can your package be downloaded using torrent?",
				},
			},
			{
				Name: "dependencies",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "[dependencies] Chiro dependencies are necessaries in your package? \n List them like this: dep1 dep2 dep3.",
				},
			},
			{
				Name: "dependencies-external",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "[dependencies] Are dependencies from other package managers necessary? \n List them like this: package-manager1:dep1 package-manager2:dep2 package-manager3|package-manager-4:dep3. Dependencies must be named using the correspondent package manager standards.",
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
			Path                 string
			Name                 string
			Version              string
			Description          string
			Url                  string
			File                 string
			Torrent              string
			Dependencies         string
			DependenciesExternal string

			Confirm bool
		}{}
		err := survey.Ask(questions, &answers)

		if err != nil {
			fmt.Println("Wizard quit because of an error: ", err)
			return
		}

		r := recipe.Recipe{
			Name:        answers.Name,
			Version:     answers.Version,
			Description: answers.Description,
			Source: recipe.Source{
				Url:     answers.Url,
				Torrent: answers.Torrent,
			},
			Security:     recipe.Security{},
			Dependencies: recipe.Dependencies{},
		}

		if answers.Dependencies != "" {
			elements := strings.SplitSeq(answers.DependenciesExternal, " ")
			for element := range elements {
				element = strings.TrimSpace(element)
				r.Dependencies.Chiro = append(r.Dependencies.Chiro, element)
			}
		}
		if answers.DependenciesExternal != "" {
			elements := strings.SplitSeq(answers.DependenciesExternal, " ")
			for element := range elements {
				element = strings.TrimSpace(element)
				r.Dependencies.System = append(r.Dependencies.System, element)
			}
		}

		if _, err := os.Stat(answers.Path); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Error while recovering the package: ", err)
			return
		}
		packData, err := os.ReadFile(answers.Path)
		if err != nil {
			fmt.Println("Error while recovering the package: ", err)
			return
		}

		privKey, _, err := filesystem.GetKeys(filesystem.GetPrivateKey)
		if err != nil {
			fmt.Println("Could not get keys, did you successfully complete setup?")
			return
		}

		data, err := crypto.SignData(privKey, packData)
		if err != nil {
			fmt.Println("Could not sign your package: " + err.Error())
			return
		}

		hash, err := crypto.EvaluateHash(packData)
		if err != nil {
			fmt.Println("Could not calculate your package's hash: " + err.Error())
			return
		}

		r.Security.Signature = data
		r.Security.SHA256 = hash

		err = recipe.WriteRecipeFile(answers.Path, &r)
		if err != nil {
			fmt.Println("Error while writing recipe file: " + err.Error())
			return
		}

		if answers.Confirm {
			fmt.Printf("\n✨ Excellent! Your package recipe has been created!\n\n")
		}
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate [dev@package#version]",
	Short: "Validate a package file",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
