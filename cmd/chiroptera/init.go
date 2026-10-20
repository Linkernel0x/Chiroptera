package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Linkernel0x/chiroptera/pkg/filesystem"
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
		}

		if force {
			confirmation := struct {
				confirm  bool
				confirm2 bool
			}{}

			questions := []*survey.Question{
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

			err := survey.Ask(questions, &confirmation)
			if err != nil {
				fmt.Println("Error while opening confirmation prompt: ", err)
				return
			}

			if !confirmation.confirm || confirmation.confirm2 {
				fmt.Println(" Aborting...")
				return
			}
		}

		err := filesystem.PurgeFolder()
		if err != nil {
			fmt.Println("Unable to purge configuration directory: ", err)
			return
		}

		
	},
}
