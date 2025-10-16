/*
Copyright © 2025 Danilo Barolo Martins de Lima <danilobml@hotmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/danilobml/lga/lga/internal/analyser"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lga",
	Short: "A CLI tool to analyse logs saved to a file",
	Long:  `lga (Log Analyser) is a command-line utility for parsing and analysing log files. 
It helps developers and system administrators identify errors, detect patterns, 
and extract useful insights from raw logs with ease.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		
		filePath := args[0]

		err := analyser.AnalyseFileLogs(filePath)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.log-analyser.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
