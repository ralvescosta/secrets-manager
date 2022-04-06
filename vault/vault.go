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

var filePath = "../example_project/.env.development"
var vaultSeparator = "$vault."
var pathKeyValueSeparator = "."
var kvVersion = 1
var vaultHost = "http://localhost:8200"
var vaultToken = "hvs.xBXPKs1UKyXqKbMCBbdWHGlg"
var fileKeyValueSeparator = "= "

type environment struct {
	vaultKey   string
	vaultValue string
	replacer   string
}

func Run() {
	envs, err := readEnvFile()
	if err != nil {
		log.Fatal(err)
	}

	if err := getVaultSecrets(envs); err != nil {
		log.Fatal(err)
	}

	if err := updateEnvFile(envs); err != nil {
		log.Fatal(err)
	}
}

func readEnvFile() (map[string][]*environment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("[secretsManager::vault] - io error: %v", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	envs := make(map[string][]*environment)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		env := strings.Split(line, vaultSeparator)
		pathAndKey := strings.Split(env[1], pathKeyValueSeparator)

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
			replacer: strings.Split(line, fileKeyValueSeparator)[1],
		})
	}
	if scanner.Err() != nil {
		return nil, fmt.Errorf("[secretsManager::vault] - file scanner error: %v", err.Error())
	}

	return envs, nil
}

func getVaultSecrets(envs map[string][]*environment) error {
	for key, values := range envs {
		res, err := getKeys(key)
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

func updateEnvFile(envs map[string][]*environment) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("[secretsManager::vault] - io error: %v", err.Error())
	}

	for _, value := range envs {
		for _, env := range value {
			file = bytes.Replace(file, []byte(env.replacer), []byte(env.vaultValue), -1)
		}
	}

	if err = ioutil.WriteFile(filePath, file, 0666); err != nil {
		return fmt.Errorf("[secretsManager::vault] - io error: %v", err.Error())
	}

	return nil
}

type vaultModel struct {
	RequestId     string            `json:"request_id"`
	LeaseId       string            `json:"lease_id"`
	LeaseDuration int               `json:"lease_duration"`
	Data          map[string]string `json:"data"`
}

func getKeys(path string) (*vaultModel, error) {
	client := &http.Client{Timeout: time.Duration(1) * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/v1/kv/%v", vaultHost, path), nil)
	if err != nil {
		return nil, fmt.Errorf("[secretsManager::vault] - internal error: %v", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Vault-Token", vaultToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[secretsManager::vault] - erro whiling get KV: %v", err.Error())
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("[secretsManager::vault] - response parser error: %v", err.Error())
	}

	body := &vaultModel{}
	err = json.Unmarshal(data, body)
	if err != nil {
		return nil, fmt.Errorf("[secretsManager::vault] - response parser error: %v", err.Error())
	}

	return body, nil
}
