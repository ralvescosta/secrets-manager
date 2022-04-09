/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package main

import (
	"log"

	"github.com/ralvescosta/secrets-manager/vault/cmd"
	"github.com/ralvescosta/secrets-manager/vault/pkg/vault"
)

func main() {
	mergeCmd := cmd.NewMergeCmd(vault.Runner)
	err := cmd.Execute(mergeCmd)

	if err != nil {
		log.Fatal(err)
	}
}
