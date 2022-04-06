/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// mergeCmd represents the base command when called without any subcommands
var mergeCmd = &cobra.Command{
	Use:   "merge", // Vault Manager
	Short: "Merge you environment file",
	Long:  `Example: cli -v https://vault -t validToken`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	mergeCmd.Flags().StringP("file-path", "f", "", "Environment file path (required)")
	mergeCmd.Flags().StringP("vault-separator", "s", "$vault.", "Vault separator pattern")
	mergeCmd.Flags().StringP("path-key-value-separator", "p", ".", "Path and Key value separator pattern")
	mergeCmd.Flags().StringP("file-key-value-separator", "q", "= ", "File key value separator pattern")
	mergeCmd.Flags().StringP("vault-host", "v", "localhost:8200", "Vault Host")
	mergeCmd.Flags().StringP("token", "t", "", "Vault Token (required)")
	mergeCmd.Flags().StringP("kv-version", "k", "1", "Key Value version")
	mergeCmd.MarkFlagRequired("file-path")
	mergeCmd.MarkFlagRequired("token")
}
