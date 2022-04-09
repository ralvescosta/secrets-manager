/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package cmd

import (
	"errors"
	"testing"

	"github.com/ralvescosta/secrets-manager/vault/pkg/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Test_NewMergeCmd_should_execute_mergeCmd_correctly(t *testing.T) {
	mergeCmd := NewMergeCmd(func(c *vault.Configs) error { return nil })
	mergeCmd.Flags().Set("file-path", "path")
	mergeCmd.Flags().Set("token", "token")

	err := mergeCmd.Execute()

	if err != nil {
		t.Error()
	}
}

func Test_getVaultConfigs_should_return_err_when_flag_missing(t *testing.T) {
	cmd := &cobra.Command{}

	configs, err := getVaultConfigs(cmd.Flags())
	if configs != nil || err == nil || err.Error() != "flag filePath is required" {
		t.Error()
	}
	cmd.Flags().String("file-path", "path", "")

	configs, err = getVaultConfigs(cmd.Flags())
	if configs != nil || err == nil || err.Error() != "wrong vault separator" {
		t.Error()
	}
	cmd.Flags().String("vault-separator", "separator", "")

	configs, err = getVaultConfigs(cmd.Flags())
	if configs != nil || err == nil || err.Error() != "wrong path key value separator" {
		t.Error()
	}
	cmd.Flags().String("path-key-value-separator", "separator", "")

	configs, err = getVaultConfigs(cmd.Flags())
	if configs != nil || err == nil || err.Error() != "wrong kv version" {
		t.Error()
	}
	cmd.Flags().String("kv-version", "version", "")

	configs, err = getVaultConfigs(cmd.Flags())
	if configs != nil || err == nil || err.Error() != "wrong vault host" {
		t.Error()
	}
	cmd.Flags().String("vault-host", "host", "")

	configs, err = getVaultConfigs(cmd.Flags())
	if configs != nil || err == nil || err.Error() != "flag token is required" {
		t.Error()
	}
	cmd.Flags().String("token", "token", "")

	configs, err = getVaultConfigs(cmd.Flags())
	if configs != nil || err == nil || err.Error() != "wrong file key value separator" {
		t.Error()
	}
}

func Test_NewMergeCmd_should_return_err_if_some_err_occur(t *testing.T) {
	mergeCmd := NewMergeCmd(func(c *vault.Configs) error { return nil })
	mergeCmd.Flags().Set("file-path", "path")
	mergeCmd.Flags().Set("token", "token")

	getVaultConfigs = func(flags *pflag.FlagSet) (*vault.Configs, error) {
		return nil, errors.New("some error")
	}

	err := mergeCmd.Execute()

	if err == nil {
		t.Error()
	}
}
