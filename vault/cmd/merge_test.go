/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package cmd

import (
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

func Test_getVaultConfigs_should_return_err_when_flag_missing(t *testing.T) {
	flags := &pflag.FlagSet{}

	_, err := getVaultConfigs(flags)
	if err == nil && err.Error() != "flag filePath is required" {
		t.Error()
	}

	flags.Set("file-path", "file")
	_, err = getVaultConfigs(flags)
	if err == nil && err.Error() != "wrong vault separator" {
		t.Error()
	}

	flags.Set("vault-separator", "separator")
	_, err = getVaultConfigs(flags)
	if err == nil && err.Error() != "wrong path key value separator" {
		t.Error()
	}

	flags.Set("path-key-value-separator", "separator")
	_, err = getVaultConfigs(flags)
	if err == nil && err.Error() != "wrong kv version" {
		t.Error()
	}

	flags.Set("key-version", "version")
	_, err = getVaultConfigs(flags)
	if err == nil && err.Error() != "wrong vault host" {
		t.Error()
	}

	flags.Set("vault-host", "host")
	_, err = getVaultConfigs(flags)
	if err == nil && err.Error() != "flag token is required" {
		t.Error()
	}

	flags.Set("token", "token")
	_, err = getVaultConfigs(flags)
	if err == nil && err.Error() != "wrong file key value separator" {
		t.Error()
	}
}
