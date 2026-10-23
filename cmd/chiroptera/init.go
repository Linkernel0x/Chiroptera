package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Linkernel0x/chiroptera/pkg/crypto"
	"github.com/Linkernel0x/chiroptera/pkg/filesystem"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.PersistentFlags().BoolVar(&force, "force", false, "Force the execution of the command even if a configuration is already existing")
}

var force bool
var setupCmd = &cobra.Command{
	Use:              "init",
	Short:            "Run chiroptera setup wizard",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🦇  Welcome to Chiroptera setup wizard!")

		if filesystem.CheckConfiguration() && !force {
			fmt.Println(" Configuration found so the command quit. To force execution of this command use the --force flag.")
			return
		}

		if force {
			confirmation := struct {
				Confirm  bool
				Confirm2 bool
			}{}

			preQuestions := []*survey.Question{
				{
					Name: "confirm",
					Prompt: &survey.Confirm{
						Renderer: survey.Renderer{},
						Message:  "This will completely erase your current configuration, deleting your keys. Are you 100% sure you want to continue?",
						Default:  false,
					},
					Validate: survey.Required,
				},
				{
					Name: "confirm2",
					Prompt: &survey.Confirm{
						Renderer: survey.Renderer{},
						Message:  "Just double checking, did you understand and want to continue? If you didn't copy your keys it will impossible to authenticate your previous identity.",
						Default:  false,
					},
					Validate: survey.Required,
				},
			}

			err := survey.Ask(preQuestions, &confirmation)
			if err != nil {
				fmt.Println("Error while opening confirmation prompt: ", err)
				return
			}

			if !confirmation.Confirm || !confirmation.Confirm2 {
				fmt.Println(" Aborting...")
				return
			}
		}

		cleaningBar := progressbar.DefaultSilent(
			-1,
			"cleaning directory",
		)

		err := filesystem.PurgeFolder()
		if err != nil {
			fmt.Println("Unable to purge configuration directory: ", err)
			return
		}

		err = cleaningBar.Finish()
		if err != nil {
			fmt.Println("Error while cleaning configuration directory: ", err)
			return
		}

		questions := []*survey.Question{
			{
				Name: "name",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "What's your dev nickname?",
				},
				Validate: survey.Required,
			},
			{
				Name: "github",
				Prompt: &survey.Input{
					Renderer: survey.Renderer{},
					Message:  "What's your GitHub account? (you're required to have one in order to authenticate)",
				},
				Validate: survey.Required,
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
			Name   string
			GitHub string

			Confirm bool
		}{}
		err = survey.Ask(questions, &answers)

		pub, priv, err := crypto.GenerateKeyPair()
		if err != nil {
			fmt.Println("Error generating key pair: ", err)
			return
		}
		err, err2 := filesystem.SaveKeys(pub, priv)
		if err != nil && err2 != nil {
			fmt.Println("Error saving keys: ", err, " ; ", err2)
			return
		} else if err2 != nil {
			fmt.Println("Error saving keys: ", err2)
		} else if err != nil {
			fmt.Println("Error saving keys: ", err)
		}

		config := filesystem.Configuration{
			SchemaVersion: 1,
			Developer:     answers.Name,
		}
		err = filesystem.SaveConfiguration(&config)
		if err != nil {
			fmt.Println("Error saving configuration: ", err)
			return
		}

		if filesystem.CheckConfiguration() {
			fmt.Println(" ✨All good! Configuration saved and ready to be used.")
		} else {
			fmt.Println("An error has occurred")
		}
		return
	},
}
