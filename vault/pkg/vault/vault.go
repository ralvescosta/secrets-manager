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

		vaultPaths := ""
		for i := 0; i < len(pathAndKey)-1; i++ {
			if i < len(pathAndKey)-2 {
				vaultPaths += pathAndKey[i] + "/"
			} else {
				vaultPaths += pathAndKey[i]
			}
		}

		envs[vaultPaths] = append(envs[vaultPaths], &environment{
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

type VaultModel struct {
	RequestId     string
	LeaseId       string
	LeaseDuration int
	Data          map[string]string
}

var getKeys = func(cfs *Configs, path string) (*VaultModel, error) {
	client := &http.Client{Timeout: time.Duration(1) * time.Second}
	req, err := http.NewRequest(http.MethodGet, getURL(cfs, path), nil)
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

	return unmarshalVaultBody(cfs.KVVersion, data)
}

var getURL = func(cfs *Configs, path string) string {
	switch cfs.KVVersion {
	case "1":
		return fmt.Sprintf("%v/v1/%v", cfs.VaultHost, path)
	case "2":
		s := strings.Split(path, "/")
		r := cfs.VaultHost + "/v1/" + s[0] + "/data"
		for i := 1; i < len(s); i++ {
			if s[i] != "" {
				r += "/" + s[i]
			}
		}
		return r
	default:
		return cfs.VaultHost
	}
}

func unmarshalVaultBody(version string, data []byte) (*VaultModel, error) {
	var err error
	switch version {
	case "1":
		bV1 := &VaultModelGen[VaultDataV1]{}
		err = json.Unmarshal(data, bV1)
		if err == nil {
			return bV1.ToModel(), nil
		}
	case "2":
		bV2 := &VaultModelGen[VaultDataV2]{}
		err = json.Unmarshal(data, bV2)
		if err == nil {
			return bV2.ToModel(), nil
		}
	}

	return nil, fmt.Errorf("response parser error: %e", err)
}

type VaultModelGen[D VaultDataGen] struct {
	RequestId     string `json:"request_id"`
	LeaseId       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Data          D      `json:"data"`
}

type VaultDataV1 map[string]string

func (pst VaultDataV1) ToModel() map[string]string {
	return pst
}

type VaultDataV2 struct {
	Data map[string]string `json:"data"`
}

func (pst VaultDataV2) ToModel() map[string]string {
	return pst.Data
}

type VaultDataGen interface {
	VaultDataV1 | VaultDataV2
	ToModel() map[string]string
}

func (pst VaultModelGen[D]) ToModel() *VaultModel {
	return &VaultModel{
		pst.RequestId,
		pst.LeaseId,
		pst.LeaseDuration,
		pst.Data.ToModel(),
	}
}
