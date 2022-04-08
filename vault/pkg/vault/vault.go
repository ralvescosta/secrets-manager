/*
Copyright © 2022 Rafael Costa <rafael.rac.mg@gmail.com>

*/
package vault

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Configs struct {
	FilePath              string
	VaultSeparator        string
	PathKeyValueSeparator string
	KVVersion             string
	VaultHost             string
	VaultToken            string
	FileKeyValueSeparator string
}

type environment struct {
	vaultKey   string
	vaultValue string
	replacer   string
}

func Runner(configs *Configs) error {
	envs, err := readEnvFile(configs)
	if err != nil {
		log.Printf("[Err] [secretesManager::vault::Runner]\n%e", err)
		return err
	}

	if err := getVaultSecrets(configs, envs); err != nil {
		log.Printf("[Err] [secretesManager::vault::Runner]\n%e", err)
		return err
	}

	if err := updateEnvFile(configs, envs); err != nil {
		log.Printf("[Err] [secretesManager::vault::Runner]\n%e", err)
		return err
	}

	return nil
}

var readEnvFile = func(cfs *Configs) (map[string][]*environment, error) {
	file, err := os.Open(cfs.FilePath)
	if err != nil {
		return nil, fmt.Errorf("io error: %v", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	envs := make(map[string][]*environment)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		env := strings.Split(line, cfs.VaultSeparator)
		pathAndKey := strings.Split(env[1], cfs.PathKeyValueSeparator)

		vaultPath := ""
		for i, v := range pathAndKey {
			if i < len(pathAndKey)-1 {
				if vaultPath == "" {
					vaultPath = v
				} else {
					vaultPath += "/" + v
				}
			}
		}

		envs[pathAndKey[0]] = append(envs[pathAndKey[0]], &environment{
			vaultKey: pathAndKey[len(pathAndKey)-1],
			replacer: strings.Split(line, cfs.FileKeyValueSeparator)[1],
		})
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("file scanner error: %v", err.Error())
	}

	return envs, nil
}

var getVaultSecrets = func(cfs *Configs, envs map[string][]*environment) error {
	for key, values := range envs {
		res, err := getKeys(cfs, key)
		if err != nil {
			return err
		}

		for _, v := range values {
			secret, ok := res.Data[v.vaultKey]
			v.vaultValue = secret
			if !ok {
				log.Printf("[secretesManager::vault] warning the secret: %v, was not find in vault", v)
			}
		}
	}

	return nil
}

var updateEnvFile = func(cfs *Configs, envs map[string][]*environment) error {
	file, err := ioutil.ReadFile(cfs.FilePath)
	if err != nil {
		return fmt.Errorf("io error: %v", err.Error())
	}

	for _, value := range envs {
		for _, env := range value {
			file = bytes.Replace(file, []byte(env.replacer), []byte(env.vaultValue), -1)
		}
	}

	if err = ioutil.WriteFile(cfs.FilePath, file, 0666); err != nil {
		return fmt.Errorf("io error: %v", err.Error())
	}

	return nil
}

type vaultModel struct {
	RequestId     string            `json:"request_id"`
	LeaseId       string            `json:"lease_id"`
	LeaseDuration int               `json:"lease_duration"`
	Data          map[string]string `json:"data"`
}

var getKeys = func(cfs *Configs, path string) (*vaultModel, error) {
	client := &http.Client{Timeout: time.Duration(1) * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/v1/kv/%v", cfs, path), nil)
	if err != nil {
		return nil, fmt.Errorf("internal error: %v", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Vault-Token", cfs.VaultToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro whiling get KV: %v", err.Error())
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("response parser error: %v", err.Error())
	}

	body := &vaultModel{}
	err = json.Unmarshal(data, body)
	if err != nil {
		return nil, fmt.Errorf("response parser error: %v", err.Error())
	}

	return body, nil
}
