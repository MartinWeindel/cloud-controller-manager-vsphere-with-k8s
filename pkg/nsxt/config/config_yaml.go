/*
 Copyright 2020 The Kubernetes Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package config

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

// validateConfig checks NSXT configurations
func (cfg *Nsxt) validateConfig() error {
	if cfg.VMCAccessToken != "" {
		if cfg.VMCAuthHost == "" {
			return errors.New("vmc auth host must be provided if auth token is provided")
		}
	} else if cfg.User != "" {
		if cfg.Password == "" {
			return errors.New("password is empty")
		}
	} else if cfg.ClientAuthKeyFile != "" {
		if cfg.ClientAuthCertFile == "" {
			return errors.New("client cert file is required if client key file is provided")
		}
	} else if cfg.ClientAuthCertFile != "" {
		if cfg.ClientAuthKeyFile == "" {
			return errors.New("client key file is required if client cert file is provided")
		}
	} else if cfg.SecretName != "" {
		if cfg.SecretNamespace == "" {
			return errors.New("secret namespace is required if secret name is provided")
		}
	} else if cfg.SecretNamespace != "" {
		if cfg.SecretName == "" {
			return errors.New("secret name is required if secret namespace is provided")
		}
	} else {
		return errors.New("user or vmc access token or client cert file must be set")
	}
	if cfg.Host == "" {
		return errors.New("host is empty")
	}
	return nil
}

// CompleteAndValidate sets default values, overrides by env and validates the resulting config
func (ncy *NsxtConfig) CompleteAndValidate() error {
	return ncy.NSXT.validateConfig()
}

// ReadRawConfig parses vSphere cloud config file and stores it into ConfigYAML
func ReadRawConfig(configData []byte) (*NsxtConfig, error) {
	if len(configData) == 0 {
		return nil, fmt.Errorf("Invalid YAML file")
	}

	cfg := NsxtConfig{}

	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		return nil, err
	}

	err := cfg.CompleteAndValidate()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ReadConfig parses vSphere cloud config file and stores it into Config
func ReadConfig(configData []byte) (*NsxtConfig, error) {
	cfg, err := ReadRawConfig(configData)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
