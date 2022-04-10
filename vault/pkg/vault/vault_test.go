/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package vault

import (
	"errors"
	"testing"
)

func Test_readEnvFile_read_file_correctly(t *testing.T) {
	_, err := readEnvFile(&Configs{
		FilePath: "../../../example/.env.development",
	})

	if err != nil {
		t.Error()
	}
}

func Test_updateEnvFile_update_file_correctly(t *testing.T) {
	err := updateEnvFile(&Configs{
		FilePath: "../../../example/.env.development",
	}, map[string][]*environment{
		"": {
			{
				vaultKey:   "",
				vaultValue: "",
				replacer:   "",
			},
		},
	})

	if err != nil {
		t.Error()
	}
}

func Test_Runner_should_execute_correctly(t *testing.T) {
	readEnvFile = func(cfs *Configs) (map[string][]*environment, error) {
		return make(map[string][]*environment), nil
	}
	getVaultSecrets = func(cfs *Configs, envs map[string][]*environment) error {
		return nil
	}
	updateEnvFile = func(cfs *Configs, envs map[string][]*environment) error {
		return nil
	}

	err := Runner(&Configs{})

	if err != nil {
		t.Error()
	}
}

func Test_Runner_should_return_err_if_some_err_occur_in_readEnvFile(t *testing.T) {
	readEnvFile = func(cfs *Configs) (map[string][]*environment, error) {
		return make(map[string][]*environment), errors.New("some error")
	}
	err := Runner(&Configs{})

	if err == nil {
		t.Error()
	}
}

func Test_Runner_should_return_err_if_some_err_occur_in_getVaultSecrets(t *testing.T) {
	readEnvFile = func(cfs *Configs) (map[string][]*environment, error) {
		return make(map[string][]*environment), nil
	}
	getVaultSecrets = func(cfs *Configs, envs map[string][]*environment) error {
		return errors.New("some error")
	}

	err := Runner(&Configs{})

	if err == nil {
		t.Error()
	}
}

func Test_Runner_should_return_err_if_some_err_occur_in_updateEnvFile(t *testing.T) {
	readEnvFile = func(cfs *Configs) (map[string][]*environment, error) {
		return make(map[string][]*environment), nil
	}
	getVaultSecrets = func(cfs *Configs, envs map[string][]*environment) error {
		return nil
	}
	updateEnvFile = func(cfs *Configs, envs map[string][]*environment) error {
		return errors.New("some error")
	}

	err := Runner(&Configs{})

	if err == nil {
		t.Error()
	}
}
