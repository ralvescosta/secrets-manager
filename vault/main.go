/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package main

import (
	"github.com/ralvescosta/secrets-manager/vault/cmd"
	"github.com/ralvescosta/secrets-manager/vault/pkg/vault"
)

func main() {
	mergeCmd := cmd.NewMergeCmd(vault.Runner)
	cmd.Execute(mergeCmd)
}
