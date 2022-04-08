/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package cmd

import (
	"errors"
	"testing"

	"github.com/ralvescosta/secrets-manager/vault/pkg/vault"
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

func Test_getVaultConfigs_should_return_err_when_flag_missing(t *testing.T) {
	flags := &pflag.FlagSet{}

	configs, err := getVaultConfigs(flags)
	if configs != nil && err == nil && err.Error() != "flag filePath is required" {
		t.Error()
	}

	flags.Set("file-path", "file")
	configs, err = getVaultConfigs(flags)
	if configs != nil && err == nil && err.Error() != "wrong vault separator" {
		t.Error()
	}

	flags.Set("vault-separator", "separator")
	configs, err = getVaultConfigs(flags)
	if configs != nil && err == nil && err.Error() != "wrong path key value separator" {
		t.Error()
	}

	flags.Set("path-key-value-separator", "separator")
	configs, err = getVaultConfigs(flags)
	if configs != nil && err == nil && err.Error() != "wrong kv version" {
		t.Error()
	}

	flags.Set("key-version", "version")
	configs, err = getVaultConfigs(flags)
	if configs != nil && err == nil && err.Error() != "wrong vault host" {
		t.Error()
	}

	flags.Set("vault-host", "host")
	configs, err = getVaultConfigs(flags)
	if configs != nil && err == nil && err.Error() != "flag token is required" {
		t.Error()
	}

	flags.Set("token", "token")
	configs, err = getVaultConfigs(flags)
	if configs != nil && err == nil && err.Error() != "wrong file key value separator" {
		t.Error()
	}
}
