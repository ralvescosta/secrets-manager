/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package cmd

import (
	"errors"
	"log"

	"github.com/ralvescosta/secrets-manager/vault/pkg/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// mergeCmd represents the base command when called without any subcommands
func NewMergeCmd(handler func(*vault.Configs) error) *cobra.Command {
	var mergeCmd = &cobra.Command{
		Use:   "merge", // Vault Manager
		Short: "Merge you environment file",
		Long:  `Example: cli -v https://vault -t validToken`,
		RunE: func(cmd *cobra.Command, args []string) error {
			configs, err := getVaultConfigs(cmd.Flags())
			if err != nil {
				log.Printf("[Err] [secretesManager::merge::Execute]\n%e", err)
				return err
			}

			return handler(configs)
		},
	}
	mergeCmd.Flags().StringP("file-path", "f", "", "Environment file path (REQUIRED)")
	mergeCmd.Flags().StringP("vault-separator", "s", "$vault.", "Vault separator pattern")
	mergeCmd.Flags().StringP("path-key-value-separator", "p", ".", "Path and Key value separator pattern")
	mergeCmd.Flags().StringP("file-key-value-separator", "q", "= ", "File key value separator pattern")
	mergeCmd.Flags().StringP("vault-host", "v", "http://localhost:8200", "Vault Host")
	mergeCmd.Flags().StringP("token", "t", "", "Vault Token (REQUIRED)")
	mergeCmd.Flags().StringP("kv-version", "k", "2", "Key Value version")

	mergeCmd.MarkFlagRequired("file-path")
	mergeCmd.MarkFlagRequired("token")

	return mergeCmd
}

var getVaultConfigs = func(flags *pflag.FlagSet) (*vault.Configs, error) {
	filePath, err := flags.GetString("file-path")
	if err != nil {
		return nil, errors.New("flag filePath is required")
	}

	vaultSeparator, err := flags.GetString("vault-separator")
	if err != nil {
		return nil, errors.New("wrong vault separator")
	}

	pathKeyValueSeparator, err := flags.GetString("path-key-value-separator")
	if err != nil {
		return nil, errors.New("wrong path key value separator")
	}

	kvVersion, err := flags.GetString("kv-version")
	if err != nil {
		return nil, errors.New("wrong kv version")
	}

	vaultHost, err := flags.GetString("vault-host")
	if err != nil {
		return nil, errors.New("wrong vault host")
	}

	vaultToken, err := flags.GetString("token")
	if err != nil {
		return nil, errors.New("flag token is required")
	}

	fileKeyValueSeparator, err := flags.GetString("file-key-value-separator")
	if err != nil {
		return nil, errors.New("wrong file key value separator")
	}

	return &vault.Configs{
		FilePath:              filePath,
		VaultSeparator:        vaultSeparator,
		PathKeyValueSeparator: pathKeyValueSeparator,
		KVVersion:             vault.MapKVVersion(kvVersion),
		VaultHost:             vaultHost,
		VaultToken:            vaultToken,
		FileKeyValueSeparator: fileKeyValueSeparator,
	}, nil
}
