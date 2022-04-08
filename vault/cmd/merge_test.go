/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package cmd

import (
	"testing"

	"github.com/ralvescosta/secrets-manager/vault/pkg/vault"
)

func Test_should_execute_mergeCmd_correctly(t *testing.T) {
	mergeCmd := NewMergeCmd(func(c *vault.Configs) error { return nil })
	mergeCmd.Flags().Set("file-path", "path")
	mergeCmd.Flags().Set("token", "token")

	err := mergeCmd.Execute()

	if err != nil {
		t.Error()
	}
}
