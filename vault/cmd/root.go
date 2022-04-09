/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "manager", // Vault Manager
	Short:   "Secrets Replacer",
	Long:    `This cli is responsable to merge your environment file with you vault secrets`,
	Version: "0.0.1",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(mergeCmd *cobra.Command) error {
	rootCmd.AddCommand(mergeCmd)
	return rootCmd.Execute()
}
