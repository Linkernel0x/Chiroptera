package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//"	  _   ,_,   _")
//"	 / `'=) (='` \\  ")
//"	/.-.-.\\ /.-.-.\\")
//" `      \"      `")

var rootCmd = &cobra.Command{
	Use:                        "chiro",
	Short:                      "A decentralized, secure, user-friendly and lightweight package manager",
	Long:                       "Chiroptera is a secure, lightweight, decentralised package manager for all kinds of developers and their users.",
	Example:                    "chiro <subcommands> [parameters] [flags]",
	SuggestionsMinimumDistance: 0,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintln(os.Stderr, err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
