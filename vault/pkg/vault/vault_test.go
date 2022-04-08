/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package vault

import "testing"

func Test_should_execute_runner_correctly(t *testing.T) {
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
