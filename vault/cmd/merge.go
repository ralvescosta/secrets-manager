/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package cmd

import (
	"log"

	"github.com/ralvescosta/secrets-manager/vault/pkg/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// mergeCmd represents the base command when called without any subcommands
var mergeCmd = &cobra.Command{
	Use:   "merge", // Vault Manager
	Short: "Merge you environment file",
	Long:  `Example: cli -v https://vault -t validToken`,
	Run: func(cmd *cobra.Command, args []string) {
		configs, err := getVaultConfigs(cmd.Flags())
		if err != nil {
			log.Fatal(err)
		}

		vault.Run(configs)
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

func getVaultConfigs(flags *pflag.FlagSet) (*vault.Configs, error) {
	filePath, err := flags.GetString("file-path")
	if err != nil {
		return nil, err
	}

	vaultSeparator, err := flags.GetString("vault-separator")
	if err != nil {
		return nil, err
	}

	pathKeyValueSeparator, err := flags.GetString("path-key-value-separator")
	if err != nil {
		return nil, err
	}

	kvVersion, err := flags.GetString("kv-version")
	if err != nil {
		return nil, err
	}

	vaultHost, err := flags.GetString("vault-host")
	if err != nil {
		return nil, err
	}

	vaultToken, err := flags.GetString("token")
	if err != nil {
		return nil, err
	}

	fileKeyValueSeparator, err := flags.GetString("file-path")
	if err != nil {
		return nil, err
	}

	return &vault.Configs{
		filePath, vaultSeparator, pathKeyValueSeparator, kvVersion, vaultHost, vaultToken, fileKeyValueSeparator,
	}, nil
}
